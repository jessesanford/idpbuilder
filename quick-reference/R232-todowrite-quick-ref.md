# 🔴🔴🔴 R232 QUICK REFERENCE - TODOWRITE OVERRIDE 🔴🔴🔴

## THE ONE RULE TO REMEMBER:
**IF TODOWRITE HAS PENDING ITEMS, YOU CANNOT STOP. PERIOD.**

## BEFORE EVERY RESPONSE ENDS:

### 1️⃣ CHECK TODOWRITE
```
Pending TODOs? → MUST EXECUTE THEM NOW
No pending TODOs? → May check other stop conditions
```

### 2️⃣ ELIMINATE "I WILL"
```
❌ "I will spawn..." → VIOLATION
✅ "Spawning now..." → COMPLIANCE
```

### 3️⃣ ACTION NOT DESCRIPTION
```
❌ "Next steps: spawn Code Reviewer" → VIOLATION  
✅ [Actually spawn Code Reviewer] → COMPLIANCE
```

## THE DEADLY SINS:

1. **Stopping with pending TODOs** = -100% INSTANT FAIL
2. **Saying "I will" without doing** = -100% INSTANT FAIL
3. **Summarizing instead of acting** = -100% INSTANT FAIL
4. **Describing work not doing it** = -100% INSTANT FAIL

## THE GOLDEN RULES:

1. **TodoWrite items are COMMANDS not suggestions**
2. **"I will" is a LIE - use "I am" and DO IT**
3. **Check TodoWrite BEFORE any stop decision**
4. **Process ALL pending items before stopping**

## QUICK CHECK:
```bash
if [ $(pending_todos_count) -gt 0 ]; then
    echo "❌ CANNOT STOP - MUST PROCESS TODOS"
    process_all_todos()
    continue_working()
else
    echo "✅ No pending TODOs - may check R021"
fi
```

## VIOLATION EXAMPLE:
```
TodoWrite: [⏳ Spawn Code Reviewer]
Orchestrator: "I will spawn Code Reviewer next"
[STOPS] 
= AUTOMATIC FAILURE
```

## CORRECT EXAMPLE:
```
TodoWrite: [⏳ Spawn Code Reviewer]
Orchestrator: "Spawning Code Reviewer now..."
[ACTUALLY SPAWNS]
= COMPLIANCE
```

---
**REMEMBER: Pending TODOs = MUST DO NOW, not later!**