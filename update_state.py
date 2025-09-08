import yaml
import sys
from datetime import datetime

# Read the current state
with open('orchestrator-state.json', 'r') as f:
    state = yaml.safe_load(f)

# Update wave_review section
if 'wave_review' not in state:
    state['wave_review'] = {}
    
state['wave_review']['report_file'] = sys.argv[1]
state['wave_review']['decision'] = sys.argv[2]
state['wave_review']['reviewed_at'] = datetime.utcnow().isoformat() + 'Z'

# Update current_review section
if 'current_review' in state:
    state['current_review']['status'] = 'APPROVED'
    state['current_review']['completed_at'] = datetime.utcnow().isoformat() + 'Z'
    state['current_review']['decision'] = sys.argv[2]

# Prepare for next state transition
state['state_machine']['next_state'] = 'WAVE_START'
state['state_machine']['next_action'] = 'Start Wave 2 implementation'
state['current_wave'] = 2  # Moving to wave 2

# Write back
with open('orchestrator-state.json', 'w') as f:
    yaml.dump(state, f, default_flow_style=False, sort_keys=False)
    
print("✅ State file updated successfully")
