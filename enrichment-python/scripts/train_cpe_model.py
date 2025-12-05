import json
import logging
import os
import sys
import random
from typing import List, Tuple
from torch.utils.data import DataLoader

# Add parent directory to path
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from app.core.config import settings

# Setup logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def train_model():
    """
    Fine-tune SentenceTransformer on CVE Description -> CPE pairs
    Uses MultipleNegativesRankingLoss (Contrastive Learning)
    """
    try:
        from sentence_transformers import SentenceTransformer, InputExample, losses, evaluation
    except ImportError:
        logger.error("sentence-transformers not installed. Cannot train.")
        return

    # 1. Load Data
    cve_file = "cve_data.json"
    if not os.path.exists(cve_file):
        logger.error("CVE data not found")
        return

    logger.info("Loading training data...")
    with open(cve_file, 'r') as f:
        cve_data = json.load(f)

    # 2. Prepare Training Examples
    # Anchor: CVE Description ("Vulnerability in Apache Tomcat...")
    # Positive: CPE Product Name ("Apache Tomcat 9.0")
    train_examples = []
    
    for cve in cve_data:
        # Extract description
        try:
            desc = cve['cve']['description']['description_data'][0]['value']
        except (KeyError, IndexError):
            continue
            
        # Extract CPEs
        cpes = extract_cpes(cve)
        
        for cpe in cpes:
            product_text = parse_cpe_to_text(cpe)
            # Create (Anchor, Positive) pair
            train_examples.append(InputExample(texts=[desc, product_text]))

    if not train_examples:
        logger.error("No valid training examples found")
        return

    logger.info(f"Prepared {len(train_examples)} training pairs")
    
    # Split Train/Val
    random.shuffle(train_examples)
    split_idx = int(len(train_examples) * 0.9)
    train_data = train_examples[:split_idx]
    
    # 3. Setup Model
    model_name = settings.model_name # 'all-MiniLM-L6-v2'
    model = SentenceTransformer(model_name)
    
    # 4. Create DataLoader
    train_dataloader = DataLoader(train_data, shuffle=True, batch_size=16)
    
    # 5. Define Loss
    # MultipleNegativesRankingLoss is great for (Anchor, Positive) pairs
    # It assumes all other positives in the batch are negatives for the current anchor
    train_loss = losses.MultipleNegativesRankingLoss(model)
    
    # 6. Train
    logger.info("Starting fine-tuning...")
    model.fit(
        train_objectives=[(train_dataloader, train_loss)],
        epochs=1,
        warmup_steps=100,
        output_path="output/fine-tuned-cpe-model",
        show_progress_bar=True
    )
    
    logger.info("Training complete. Model saved to output/fine-tuned-cpe-model")

def extract_cpes(cve_data):
    """Helper to extract CPEs"""
    cpes = set()
    try:
        nodes = cve_data.get('configurations', {}).get('nodes', [])
        for node in nodes:
            for match in node.get('cpe_match', []):
                if match.get('vulnerable'):
                    cpes.add(match.get('cpe23Uri'))
    except Exception:
        pass
    return list(cpes)

def parse_cpe_to_text(cpe: str) -> str:
    """Helper to convert CPE to text"""
    parts = cpe.split(':')
    if len(parts) >= 6:
        vendor = parts[3].replace('_', ' ')
        product = parts[4].replace('_', ' ')
        return f"{vendor} {product}".strip()
    return cpe

if __name__ == "__main__":
    train_model()

