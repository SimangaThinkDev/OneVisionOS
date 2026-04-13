package nids

// DefaultSignatures returns a set of basic security signatures for educational environments
func DefaultSignatures() []Signature {
	return []Signature{
		{
			ID:             "OV-001",
			Name:           "Reverse Shell Attempt",
			Description:    "Detection of common reverse shell payload patterns",
			Level:          LevelCritical,
			PayloadPattern: "/bin/bash -i",
		},
		{
			ID:             "OV-002",
			Name:           "SQL Injection Probe",
			Description:    "Detection of basic SQL injection patterns in network traffic",
			Level:          LevelHigh,
			PayloadPattern: "UNION SELECT",
		},
		{
			ID:             "OV-003",
			Name:           "SSH Brute Force Signature",
			Description:    "Unusual SSH traffic patterns (Simplified)",
			Level:          LevelMedium,
			DestPort:       22,
			PayloadPattern: "SSH-2.0-LibSSH",
		},
	}
}
