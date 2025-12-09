-- 007_seed_config_standards.sql
-- Seed data for CIS Benchmark standards for Cisco ASA, IOS, Palo Alto, and Fortinet

BEGIN;

-- ============================================================================
-- CIS CISCO ASA BENCHMARK
-- ============================================================================

-- Authentication and Authorization
INSERT INTO config_standards (
    standard_name, standard_version, manufacturer, device_type, model_family,
    category, requirement_id, requirement_title, requirement_description,
    compliance_frameworks, compliance_requirement,
    check_type, check_config_path, check_pattern, expected_value,
    default_severity, priority, remediation_guidance, remediation_example,
    status
) VALUES
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'authentication', 'CIS-ASA-1.1', 'Ensure password encryption is enabled',
    'Passwords should be encrypted using the enable password encryption command',
    '["CIS", "PCI-DSS", "NIST"]'::jsonb,
    'PCI-DSS 8.2.1: Protect passwords with strong cryptography',
    'presence', 'crypto.config', 'enable password encryption', 'enable password encryption',
    'high', 'high',
    'Enable password encryption: enable password encryption',
    'enable password encryption',
    'active'
),
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'authentication', 'CIS-ASA-1.2', 'Ensure strong password policy is configured',
    'Configure minimum password length and complexity requirements',
    '["CIS", "PCI-DSS", "NIST"]'::jsonb,
    'PCI-DSS 8.2.3: Enforce strong passwords',
    'pattern_match', 'password policy', 'password.*length.*8|password.*min.*8', 'password.*length.*8',
    'high', 'high',
    'Configure password policy: password-policy min-length 8',
    'password-policy min-length 8',
    'active'
),
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'authentication', 'CIS-ASA-1.3', 'Ensure default accounts are removed or disabled',
    'Remove or disable default user accounts (admin, cisco, etc.)',
    '["CIS", "PCI-DSS"]'::jsonb,
    'CIS Control 5: Account Management',
    'absence', 'user_accounts', 'username (admin|cisco|user|guest)', '',
    'critical', 'critical',
    'Remove default accounts: no username admin',
    'no username admin',
    'active'
),
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'authentication', 'CIS-ASA-1.4', 'Ensure local user accounts use encrypted passwords',
    'All local user passwords should be encrypted',
    '["CIS", "PCI-DSS"]'::jsonb,
    'PCI-DSS 8.2.1: Protect passwords',
    'pattern_match', 'passwords', 'password.*encrypted|password.*secret', 'password.*encrypted',
    'high', 'high',
    'Use encrypted passwords: username test password encrypted <hash>',
    'username test password encrypted <hash>',
    'active'
);

-- Network Security
INSERT INTO config_standards (
    standard_name, standard_version, manufacturer, device_type, model_family,
    category, requirement_id, requirement_title, requirement_description,
    compliance_frameworks, compliance_requirement,
    check_type, check_config_path, check_pattern, expected_value,
    default_severity, priority, remediation_guidance, remediation_example,
    status
) VALUES
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'network', 'CIS-ASA-2.1', 'Ensure Telnet is disabled',
    'Telnet transmits data in plaintext and should be disabled',
    '["CIS", "PCI-DSS", "NIST"]'::jsonb,
    'PCI-DSS 4.1: Use strong cryptography',
    'absence', 'telnet', 'telnet', '',
    'high', 'high',
    'Disable Telnet: no telnet',
    'no telnet',
    'active'
),
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'network', 'CIS-ASA-2.2', 'Ensure SSH is enabled and configured securely',
    'SSH should be enabled with strong encryption and key exchange',
    '["CIS", "PCI-DSS", "NIST"]'::jsonb,
    'PCI-DSS 4.1: Use strong cryptography',
    'presence', 'ssh', 'ssh.*enable|ssh.*version.*2', 'ssh version 2',
    'high', 'high',
    'Enable SSH: ssh version 2',
    'ssh version 2',
    'active'
),
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'network', 'CIS-ASA-2.3', 'Ensure HTTP server is disabled',
    'HTTP server should be disabled to prevent unencrypted management access',
    '["CIS", "PCI-DSS"]'::jsonb,
    'CIS Control 9: Limitation and Control of Network Ports',
    'absence', 'http_server', 'http server enable', '',
    'medium', 'medium',
    'Disable HTTP server: no http server enable',
    'no http server enable',
    'active'
),
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'network', 'CIS-ASA-2.4', 'Ensure SNMP is configured securely',
    'SNMP should use SNMPv3 with authentication and encryption',
    '["CIS", "NIST"]'::jsonb,
    'NIST AC-17: Remote Access',
    'pattern_match', 'snmp', 'snmp-server.*v3|snmp-server.*auth', 'snmp-server.*v3',
    'medium', 'medium',
    'Use SNMPv3: snmp-server host <host> version 3 auth <community>',
    'snmp-server host <host> version 3 auth <community>',
    'active'
);

-- Logging and Monitoring
INSERT INTO config_standards (
    standard_name, standard_version, manufacturer, device_type, model_family,
    category, requirement_id, requirement_title, requirement_description,
    compliance_frameworks, compliance_requirement,
    check_type, check_config_path, check_pattern, expected_value,
    default_severity, priority, remediation_guidance, remediation_example,
    status
) VALUES
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'logging', 'CIS-ASA-3.1', 'Ensure logging is enabled',
    'Enable logging to track security events and configuration changes',
    '["CIS", "PCI-DSS", "NIST", "ISO27001"]'::jsonb,
    'PCI-DSS 10.2: Implement automated audit trails',
    'presence', 'logging', 'logging enable|logging.*on', 'logging enable',
    'high', 'high',
    'Enable logging: logging enable',
    'logging enable',
    'active'
),
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'logging', 'CIS-ASA-3.2', 'Ensure syslog is configured',
    'Configure syslog server for centralized logging',
    '["CIS", "PCI-DSS", "NIST"]'::jsonb,
    'PCI-DSS 10.5: Secure audit trails',
    'presence', 'logging', 'logging.*host|logging.*server', 'logging host',
    'medium', 'medium',
    'Configure syslog: logging host <interface> <ip>',
    'logging host inside 192.168.1.100',
    'active'
);

-- Access Control
INSERT INTO config_standards (
    standard_name, standard_version, manufacturer, device_type, model_family,
    category, requirement_id, requirement_title, requirement_description,
    compliance_frameworks, compliance_requirement,
    check_type, check_config_path, check_pattern, expected_value,
    default_severity, priority, remediation_guidance, remediation_example,
    status
) VALUES
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'access_control', 'CIS-ASA-4.1', 'Ensure access-lists are configured',
    'Configure access control lists to restrict network traffic',
    '["CIS", "PCI-DSS", "NIST"]'::jsonb,
    'PCI-DSS 1.2: Build firewall configuration',
    'presence', 'access_lists', 'access-list', 'access-list',
    'high', 'high',
    'Configure access-lists: access-list <name> extended permit/deny <rule>',
    'access-list OUTSIDE extended permit tcp any any eq 443',
    'active'
),
(
    'CIS Cisco ASA Benchmark', '1.0', 'Cisco', 'firewall', 'ASA',
    'access_control', 'CIS-ASA-4.2', 'Ensure default deny rule exists',
    'Access-lists should have an implicit deny at the end',
    '["CIS", "PCI-DSS"]'::jsonb,
    'CIS Control 9: Limitation and Control of Network Ports',
    'pattern_match', 'access_lists', 'access-list.*deny.*any|implicit.*deny', 'implicit deny',
    'high', 'high',
    'Ensure implicit deny: Access-lists have implicit deny by default',
    'Access-lists have implicit deny by default',
    'active'
);

-- ============================================================================
-- CIS CISCO IOS BENCHMARK
-- ============================================================================

INSERT INTO config_standards (
    standard_name, standard_version, manufacturer, device_type, model_family,
    category, requirement_id, requirement_title, requirement_description,
    compliance_frameworks, compliance_requirement,
    check_type, check_config_path, check_pattern, expected_value,
    default_severity, priority, remediation_guidance, remediation_example,
    status
) VALUES
(
    'CIS Cisco IOS Benchmark', '1.0', 'Cisco', 'router', 'IOS',
    'authentication', 'CIS-IOS-1.1', 'Ensure password encryption is enabled',
    'Enable password encryption: service password-encryption',
    '["CIS", "PCI-DSS"]'::jsonb,
    'PCI-DSS 8.2.1: Protect passwords',
    'presence', 'authentication', 'service password-encryption', 'service password-encryption',
    'high', 'high',
    'Enable password encryption: service password-encryption',
    'service password-encryption',
    'active'
),
(
    'CIS Cisco IOS Benchmark', '1.0', 'Cisco', 'router', 'IOS',
    'network', 'CIS-IOS-2.1', 'Ensure Telnet is disabled',
    'Disable Telnet: no line vty 0 4 transport input telnet',
    '["CIS", "PCI-DSS"]'::jsonb,
    'PCI-DSS 4.1: Use strong cryptography',
    'absence', 'telnet', 'transport input telnet', '',
    'high', 'high',
    'Disable Telnet: no line vty 0 4 transport input telnet',
    'no line vty 0 4 transport input telnet',
    'active'
),
(
    'CIS Cisco IOS Benchmark', '1.0', 'Cisco', 'router', 'IOS',
    'network', 'CIS-IOS-2.2', 'Ensure SSH is enabled',
    'Enable SSH: line vty 0 4 transport input ssh',
    '["CIS", "PCI-DSS"]'::jsonb,
    'PCI-DSS 4.1: Use strong cryptography',
    'presence', 'ssh', 'transport input ssh', 'transport input ssh',
    'high', 'high',
    'Enable SSH: line vty 0 4 transport input ssh',
    'line vty 0 4 transport input ssh',
    'active'
);

-- ============================================================================
-- CIS PALO ALTO NETWORKS BENCHMARK
-- ============================================================================

INSERT INTO config_standards (
    standard_name, standard_version, manufacturer, device_type, model_family,
    category, requirement_id, requirement_title, requirement_description,
    compliance_frameworks, compliance_requirement,
    check_type, check_config_path, check_pattern, expected_value,
    default_severity, priority, remediation_guidance, remediation_example,
    status
) VALUES
(
    'CIS Palo Alto Networks Benchmark', '1.0', 'Palo Alto', 'firewall', 'PAN-OS',
    'authentication', 'CIS-PAN-1.1', 'Ensure strong password policy is configured',
    'Configure minimum password length and complexity',
    '["CIS", "PCI-DSS"]'::jsonb,
    'PCI-DSS 8.2.3: Enforce strong passwords',
    'pattern_match', 'password policy', 'password.*length|password.*complexity', 'password.*length.*8',
    'high', 'high',
    'Configure password policy in Device > Server Profiles > LDAP/Active Directory',
    'Set minimum password length to 8 characters',
    'active'
),
(
    'CIS Palo Alto Networks Benchmark', '1.0', 'Palo Alto', 'firewall', 'PAN-OS',
    'access_control', 'CIS-PAN-2.1', 'Ensure security policies are configured',
    'Configure security policies to control traffic flow',
    '["CIS", "PCI-DSS"]'::jsonb,
    'PCI-DSS 1.2: Build firewall configuration',
    'presence', 'security_policies', 'security.*policy|policy.*rule', 'security policy',
    'high', 'high',
    'Configure security policies in Policies > Security',
    'Create security policy rules',
    'active'
);

-- ============================================================================
-- CIS FORTINET FORTIGATE BENCHMARK
-- ============================================================================

INSERT INTO config_standards (
    standard_name, standard_version, manufacturer, device_type, model_family,
    category, requirement_id, requirement_title, requirement_description,
    compliance_frameworks, compliance_requirement,
    check_type, check_config_path, check_pattern, expected_value,
    default_severity, priority, remediation_guidance, remediation_example,
    status
) VALUES
(
    'CIS Fortinet FortiGate Benchmark', '1.0', 'Fortinet', 'firewall', 'FortiOS',
    'authentication', 'CIS-FGT-1.1', 'Ensure strong password policy is configured',
    'Configure minimum password length and complexity',
    '["CIS", "PCI-DSS"]'::jsonb,
    'PCI-DSS 8.2.3: Enforce strong passwords',
    'pattern_match', 'password policy', 'set.*password.*min|password.*length', 'set password-policy min-length 8',
    'high', 'high',
    'Configure password policy: config user setting, set password-policy min-length 8',
    'config user setting\nset password-policy min-length 8',
    'active'
),
(
    'CIS Fortinet FortiGate Benchmark', '1.0', 'Fortinet', 'firewall', 'FortiOS',
    'access_control', 'CIS-FGT-2.1', 'Ensure firewall policies are configured',
    'Configure firewall policies to control traffic',
    '["CIS", "PCI-DSS"]'::jsonb,
    'PCI-DSS 1.2: Build firewall configuration',
    'presence', 'firewall_policies', 'config firewall policy', 'config firewall policy',
    'high', 'high',
    'Configure firewall policies: config firewall policy',
    'config firewall policy',
    'active'
);

COMMIT;

