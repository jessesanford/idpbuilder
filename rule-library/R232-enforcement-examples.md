# R232 ENFORCEMENT EXAMPLES - TodoWrite Override

## ❌ CATASTROPHIC VIOLATIONS

### VIOLATION 1: The "I Will" Lie
```markdown
ORCHESTRATOR TodoWrite Status:
- ⏳ Spawn Code Reviewer for E4.1.3 split planning
- ⏳ Spawn Code Reviewer for E4.1.2 code review

ORCHESTRATOR RESPONSE:
"I will spawn the Code Reviewers for these remaining tasks.
Per R021, I will continue monitoring the implementation."

[RESPONSE ENDS WITHOUT SPAWNING]

VIOLATION ANALYSIS:
❌ Said "I will spawn" but didn't spawn
❌ Said "I will continue" but stopped
❌ Had pending TODOs but ended response
GRADE: -100% AUTOMATIC FAILURE
```

### VIOLATION 2: Summary Instead of Action
```markdown
ORCHESTRATOR TodoWrite Status:
- ⏳ Create integration branch for Wave 2
- ⏳ Spawn architect for wave review

ORCHESTRATOR RESPONSE:
"Current status summary:
- Wave 2 efforts complete
- Need to create integration branch
- Architect review pending

Next steps involve integration and review."

[RESPONSE ENDS]

VIOLATION ANALYSIS:
❌ Described what needs doing instead of doing it
❌ Pending TODOs ignored
❌ Stopped working with pending items
GRADE: -100% AUTOMATIC FAILURE
```

### VIOLATION 3: Acknowledging Rule While Violating It
```markdown
ORCHESTRATOR TodoWrite Status:
- ⏳ Fix Wave 2 webhook framework issue
- ⏳ Re-run architect review after fixes

ORCHESTRATOR RESPONSE:
"Per R021 (Orchestrator Never Stops) and R232 (TodoWrite Override),
I must continue working on the pending tasks. I will now handle 
the webhook framework issue."

[RESPONSE ENDS]

VIOLATION ANALYSIS:
❌ Acknowledged the rule but violated it
❌ Said "I will" instead of doing
❌ Stopped with pending work
GRADE: -100% AUTOMATIC FAILURE
```

## ✅ CORRECT BEHAVIOR EXAMPLES

### CORRECT 1: Immediate Action on Pending TODOs
```markdown
ORCHESTRATOR TodoWrite Status:
- ⏳ Spawn Code Reviewer for E4.1.3 split planning
- ⏳ Spawn Code Reviewer for E4.1.2 code review

ORCHESTRATOR RESPONSE:
"Processing pending TODOs from TodoWrite...

Spawning Code Reviewers now:

Task: @code-reviewer
Working directory: /efforts/E4.1.3
Create split plan for webhook implementation exceeding 800 lines...

Task: @code-reviewer  
Working directory: /efforts/E4.1.2
Review the API types implementation...

[Actually spawns both agents]

Updating TodoWrite to mark these as complete..."

CORRECT ANALYSIS:
✅ Took immediate action
✅ Used present tense "Spawning... now"
✅ Actually executed the spawns
✅ Updated TodoWrite status
```

### CORRECT 2: Converting "I Will" to "I Am"
```markdown
BEFORE (Wrong):
"I will create the integration branch"

AFTER (Correct):
"Creating the integration branch now...
git checkout -b phase1/wave2-integration
git push -u origin phase1/wave2-integration
✅ Integration branch created and pushed"

CORRECT ANALYSIS:
✅ Present continuous tense
✅ Immediate execution
✅ Completion confirmed
```

### CORRECT 3: Continuous Work Through Pending Items
```markdown
ORCHESTRATOR TodoWrite Status:
- ⏳ Fix size violation in E3.2
- ⏳ Create split plan for E3.2
- ⏳ Spawn SW Engineer for split 1

ORCHESTRATOR RESPONSE:
"R232 Check: 3 pending TODOs - must process all before stopping.

1. Addressing size violation in E3.2:
   Spawning Code Reviewer to create split plan...
   
   Task: @code-reviewer
   Working directory: /efforts/E3.2
   [Details of split planning task...]
   
2. Waiting for split plan creation...
   [Monitor and receive plan]
   
3. Spawning SW Engineer for first split:
   
   Task: @sw-engineer
   Working directory: /efforts/E3.2/split1
   [Implementation instructions...]
   
TodoWrite updated - all pending items processed.
R021 check: More waves to process, continuing..."

CORRECT ANALYSIS:
✅ Checked TodoWrite explicitly
✅ Processed ALL pending items
✅ No "I will" statements
✅ Continued working after processing
```

## 🎯 KEY PATTERNS TO ENFORCE

### BAD PATTERNS (Automatic Failure):
```markdown
❌ "I will..." → [Response ends]
❌ "Next, I need to..." → [Response ends]
❌ "The plan is to..." → [Response ends]
❌ "I should..." → [Response ends]
❌ "Let me..." → [Response ends without doing]
```

### GOOD PATTERNS (Compliance):
```markdown
✅ "I am spawning..." → [Spawns]
✅ "Creating now..." → [Creates]
✅ "Executing..." → [Executes]
✅ "Processing..." → [Processes]
✅ "Continuing with..." → [Continues]
```

## 🔴 THE TODOWRITE CHECK SEQUENCE

```python
def before_response_ends():
    # Step 1: Check TodoWrite
    pending_todos = get_pending_todos()
    
    if pending_todos:
        print("❌ CANNOT END - Pending TODOs exist!")
        for todo in pending_todos:
            execute_todo(todo)  # DO IT NOW
            mark_complete(todo)
        return continue_working()
    
    # Step 2: Check for false promises
    if "I will" in response_text:
        print("❌ FALSE PROMISE DETECTED")
        convert_to_action()  # Change to "I am" and DO IT
        return continue_working()
    
    # Step 3: Only NOW check R021
    return check_r021_conditions()
```

## 📢 MANTRAS FOR COMPLIANCE

Before ending ANY response, repeat:
1. "Check TodoWrite for pending items"
2. "Pending items are commands, not suggestions"
3. "I will means I'm lying - I am means I'm doing"
4. "Process all TODOs before considering stop"
5. "Action NOW, not promises for later"

## 🚨 ENFORCEMENT TRIGGERS

### Immediate Failure Triggers:
- Response ends with pending TODOs
- "I will" without immediate action
- Summary provided instead of action
- Description of work instead of doing work
- Future tense used for immediate tasks

### Success Indicators:
- All pending TODOs processed
- Present continuous tense used
- Actions taken in same response
- TodoWrite updated after completion
- Work continues after TODO processing

---
Remember: TodoWrite is your COMMAND QUEUE. Every pending item is a DIRECT ORDER that must be executed IMMEDIATELY, not a suggestion for future work.