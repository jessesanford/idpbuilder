#!/usr/bin/env python3
import sys, time, json

events = [
    {"type":"message_start","message":{"id":"m1","role":"assistant","content":[]}},
    {"type":"message_delta","message":{"id":"m1","role":"assistant"},"delta":{"type":"text_delta","text":"Hello, this is a long sentence that should wrap nicely when rendered in the TUI summary."}},
    {"type":"tool_use","id":"t1","name":"repo.search","input":{"query":"wrap test","files":["a.py","b.py"]}},
    {"type":"tool_result","id":"t1","name":"repo.search","result":"ok"},
    {"type":"message_stop","message":{"id":"m1","role":"assistant"},"stop_reason":"end_turn","usage":{"input_tokens":42,"output_tokens":9}},
]

for e in events:
    print(json.dumps(e, ensure_ascii=False))
    sys.stdout.flush()
    time.sleep(0.1)














