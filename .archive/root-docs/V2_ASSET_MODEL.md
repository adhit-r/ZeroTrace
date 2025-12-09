# V2 Asset Data Model

This document outlines the enhanced data model for assets within ZeroTrace, designed to provide a comprehensive and detailed view of each monitored system.

## 1. Core Asset Information

This is the primary object representing a physical or virtual machine.

| Field             | Type      | Description                                                 | Example                               |
| ----------------- | --------- | ----------------------------------------------------------- | ------------------------------------- |
| `id`              | `UUID`    | Unique identifier for the asset (agent ID)                  | `1df24d47-dd58-416a-8f43-108f8b438cda` |
| `organization_id` | `UUID`    | The organization the asset belongs to                       | `a9582235-313a-4e78-a28f-72f6946049a3` |
| `hostname`        | `String`  | The hostname of the machine                                 | `Adhis-MacBook-Pro.local`             |
| `ip_address`      | `String`  | The primary IP address of the asset                         | `192.168.1.10`                        |
| `mac_address`     | `String`  | The primary MAC address                                     | `A8:8B:8B:8B:8B:8B`                   |
| `tags`            | `Array`   | User-defined tags for grouping and filtering                | `["critical", "production", "pci"]`   |
| `risk_score`      | `Float`   | Calculated overall risk score for the asset                 | `8.5`                                 |
| `last_seen`       | `ISO 8601`| The last time the agent reported to the API                 | `2025-10-12T10:45:00Z`                |
| `created_at`      | `ISO 8601`| When the asset was first registered                         | `2025-09-01T12:00:00Z`                |

## 2. Operating System Details

| Field             | Type      | Description                                                 | Example                               |
| ----------------- | --------- | ----------------------------------------------------------- | ------------------------------------- |
| `os_name`         | `String`  | The name of the operating system                            | `macOS`                               |
| `os_version`      | `String`  | The specific version of the OS                              | `14.5`                                |
| `os_build`        | `String`  | The build number of the OS                                  | `23F79`                               |
| `kernel_version`  | `String`  | The version of the underlying OS kernel                     | `Darwin 23.5.0`                       |

## 3. Hardware Specifications

| Field             | Type      | Description                                                 | Example                               |
| ----------------- | --------- | ----------------------------------------------------------- | ------------------------------------- |
| `cpu_model`       | `String`  | The model name of the CPU                                   | `Apple M2 Pro`                        |
| `cpu_cores`       | `Integer` | The number of CPU cores                                     | `12`                                  |
| `memory_total_gb` | `Float`   | Total physical RAM in gigabytes                             | `16.0`                                |
| `storage_total_gb`| `Float`   | Total storage capacity in gigabytes                         | `512.1`                               |
| `gpu_model`       | `String`  | The model name of the primary GPU                           | `Apple M2 Pro GPU`                    |
| `serial_number`   | `String`  | The hardware serial number                                  | `C02G8R2JML7H`                        |
| `platform`        | `String`  | The hardware platform (e.g., `arm64`, `x86_64`)             | `arm64`                               |

## 4. Location Information (Best Effort)

Location data will be gathered via IP geolocation and may not always be precise.

| Field             | Type      | Description                                                 | Example                               |
| ----------------- | --------- | ----------------------------------------------------------- | ------------------------------------- |
| `city`            | `String`  | Estimated city based on public IP                           | `San Francisco`                       |
| `region`          | `String`  | Estimated region or state                                   | `California`                          |
| `country`         | `String`  | Estimated country                                           | `United States`                       |
| `timezone`        | `String`  | The system's configured timezone                            | `America/Los_Angeles`                 |

## 5. Associated Data (Linked by `asset_id`)

These will be stored in separate tables or documents and linked back to the core asset.

*   **Installed Applications**: A list of all discovered software, including name, version, and installation path.
*   **Vulnerabilities**: A list of all detected CVEs and configuration issues.
*   **Network Neighbors**: A list of other assets detected on the same network segment.
*   **SBOMs**: For supported applications, a detailed Software Bill of Materials.

This enhanced model will provide the necessary foundation for the detailed asset views, improved risk scoring, and advanced filtering required for a comprehensive vulnerability management platform.

