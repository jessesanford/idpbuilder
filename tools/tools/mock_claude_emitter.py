#!/usr/bin/env python3
import sys
import time
import json
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("--emit-continuation", action="store_true")
parser.add_argument("--interval", type=float, default=0.2)
parser.add_argument("--mode", choices=["simple","sf2","sf2_exact"], default="simple")
args = parser.parse_args()

if args.mode == "sf2":
    # Emit messages closely matching the Software Factory sample
    system_msg = {
        "type": "system",
        "subtype": "init",
        "cwd": "/home/vscode/workspaces/idpbuilder-push",
        "session_id": "3c350c77-fa41-4065-a987-bd07068b5127",
        "tools": [
            "Task","Bash","Glob","Grep","ExitPlanMode","Read","Edit","MultiEdit",
            "Write","NotebookEdit","WebFetch","TodoWrite","WebSearch","BashOutput","KillShell"
        ],
        "mcp_servers": [],
        "model": "claude-opus-4-1-20250805",
        "permissionMode": "bypassPermissions",
        "slash_commands": [
            "continue-orchestrating","continue-reviewing","cascade-orchestration","init-software-factory",
            "integrate-branches","check-status","pr-ready-transform","continue-architecting","integration",
            "COMMANDS-QUICK-REFERENCE","continue-implementing","reset-state","splitting","fix-cascade","README"
        ],
        "apiKeySource": "none",
        "output_style": "default",
        "uuid": "1631cb62-5aef-463f-9fdd-59b804fa7bd5",
    }
    print(json.dumps(system_msg), flush=True)

    thinking_text = (
        "The user is invoking the /continue-orchestrating command. Let me analyze the current state: "
        "Pre-flight checks and state verification will commence."
    )
    assistant_thinking = {
        "type": "assistant",
        "message": {
            "id": "msg_01Lk3kuJspdcANLmH2jiXSpK",
            "type": "message",
            "role": "assistant",
            "model": "claude-opus-4-1-20250805",
            "content": [{"type": "thinking", "thinking": thinking_text}],
        },
        "usage": {"input_tokens": 10, "output_tokens": 4},
    }
    print(json.dumps(assistant_thinking), flush=True)

    assistant_text = {
        "type": "assistant",
        "message": {
            "id": "msg_01Lk3kuJspdcANLmH2jiXSpK",
            "type": "message",
            "role": "assistant",
            "model": "claude-opus-4-1-20250805",
            "content": [{"type": "text", "text": "I'll begin the orchestrator continuation process with all mandatory checks and acknowledgments."}],
        },
    }
    print(json.dumps(assistant_text), flush=True)

    tool_use_id = "toolu_mock_bash_1"
    assistant_tool_use = {
        "type": "assistant",
        "message": {
            "id": "msg_01Lk3kuJspdcANLmH2jiXSpK",
            "type": "message",
            "role": "assistant",
            "model": "claude-opus-4-1-20250805",
            "content": [
                {
                    "type": "tool_use",
                    "id": tool_use_id,
                    "name": "Bash",
                    "input": {
                        "command": (
                            "echo \"╔═══════════════════════════════════════════════════════════════════════════════╗\n"
                            "║                        SOFTWARE FACTORY 2.0                                  ║\n"
                            "║                   ORCHESTRATOR CONTINUATION                                   ║\n"
                            "╚═══════════════════════════════════════════════════════════════════════════════╝\n\n"
                            "🎯 AGENT IDENTITY: ORCHESTRATOR\n"
                            "✓ Confirming identity: orchestrator\n\" && date '+%Y-%m-%d %H:%M:%S %Z'"
                        ),
                        "description": "Start orchestrator continuation and show timestamp",
                    },
                }
            ],
        },
    }
    print(json.dumps(assistant_tool_use), flush=True)

    tool_result = (
        "╔═══════════════════════════════════════════════════════════════════════════════╗\n"
        "║                        SOFTWARE FACTORY 2.0                                  ║\n"
        "║                   ORCHESTRATOR CONTINUATION                                   ║\n"
        "╚═══════════════════════════════════════════════════════════════════════════════╝\n\n"
        "🎯 AGENT IDENTITY: ORCHESTRATOR\n"
        "✓ Confirming identity: orchestrator\n\n"
        "2025-09-23 16:22:45 UTC"
    )
    user_tool_result = {
        "type": "user",
        "message": {
            "role": "user",
            "content": [{"tool_use_id": tool_use_id, "type": "tool_result", "content": tool_result, "is_error": False}],
        },
    }
    print(json.dumps(user_tool_result), flush=True)

    if args.emit_continuation:
        time.sleep(args.interval)
        print("CONTINUE-SOFTWARE-FACTORY=TRUE", flush=True)
elif args.mode == "sf2_exact":
    # Emit the five JSON lines exactly as provided in user's sample (one per line)
    l1 = r"""{"type":"system","subtype":"init","cwd":"/home/vscode/workspaces/idpbuilder-push","session_id":"3c350c77-fa41-4065-a987-bd07068b5127","tools":["Task","Bash","Glob","Grep","ExitPlanMode","Read","Edit","MultiEdit","Write","NotebookEdit","WebFetch","TodoWrite","WebSearch","BashOutput","KillShell"],"mcp_servers":[],"model":"claude-opus-4-1-20250805","permissionMode":"bypassPermissions","slash_commands":["continue-orchestrating","continue-reviewing","cascade-orchestration","init-software-factory","integrate-branches","check-status","pr-ready-transform","continue-architecting","integration","COMMANDS-QUICK-REFERENCE","continue-implementing","reset-state","splitting","fix-cascade","README","continue-orchestrating","COMMANDS-QUICK-REFERENCE","continue-implementing","reset-state","splitting","fix-cascade","README","pr-ready-transform","continue-architecting","integration","init-software-factory","check-status","continue-reviewing","archived:spawn-integration-agent-p1w2-cascade","archived:integrate-branches","archived:spawn-code-reviewer-p2w2-merge-plan","archived:fix-cascade-orchestration","archived:spawn-code-reviewer-update-p1w2-merge-plan","archived:bug-fix-cascade-orchestration","archived:spawn-integration-agent-p2w1-cascade","archived:sub-orchestrate-fix","archived:cascade-orchestration","archived:sub-orchestrate-integration","archived:sub-orchestrate-split","compact","context","cost","init","output-style:new","pr-comments","release-notes","todos","review","security-review"],"apiKeySource":"none","output_style":"default","uuid":"1631cb62-5aef-463f-9fdd-59b804fa7bd5"}"""
    l2 = r"""{"type":"assistant","message":{"id":"msg_01Lk3kuJspdcANLmH2jiXSpK","type":"message","role":"assistant","model":"claude-opus-4-1-20250805","content":[{"type":"thinking","thinking":"The user is invoking the /continue-orchestrating command. Let me analyze the current state:\n\n1. From the git status, I can see we're on the `software-factory-2.0` branch with many modified and untracked files\n2. There's an orchestrator-state-v3.json file that's been modified \n3. There are TODO files in the todos/ directory\n4. Many utilities and rule files have been modified\n\nFollowing the instructions, I need to:\n1. Confirm my identity as the orchestrator\n2. Perform all mandatory pre-flight checks\n3. Check for the master implementation plan\n4. Load the current state from orchestrator-state-v3.json\n5. Continue orchestration based on the current state\n\nLet me start with the pre-flight checks and state verification.","signature":"EvYGCkYIBxgCKkCBZndL1JVES8EvwuwEoYdY0cLc5aW8r7TDl/HVWiIHC3oow5FCAAnrdfJcaL2PlzbKoPieUBAb/1a2yQaCkuUxEgwQa0gAK7MmUFC1r8UaDJlOmCch80f9z55XviIwOhecXnLN2JpeebwW1zsuyYYIImeZ7GuzOwF0tf7Lsa7k2QszF/6W4kDtZ4byMtraKt0FxKBnIJXf244UfZQ7qjhiueJtCjxXhOZ22OKJ4+Q70p3OHGkg3NCx1hsGv6nbrSbxROkTblZsQKVTiqiMRt4Pzj0J8m41y4OdtDgUDMp/wi4ShP/H2CnsqS51/F1dNEm2Gz9ySnXH+wj1VcwelzKcGjo1bvIV9n05jIsiibnRsjiP9TErhwvS388+aQBvVNrDufWmPrTLiElDdmq/1HmcgoHWEf4P8W6xIqq2VsKkA+LX6Ekb2X728SDfarm38WnchbdJvsrfcRwsQr5Xbc6013YHOmdd28R+nRSl2b4PrnEtL2rbBt1le38S5JxmolcvSOmQeg/5GefGbvYuXf9JLybSu5saH7G9/j8oqj/JEnJ422eoL1MOmYOz226ZXCvgmq3IbX50NM1+ScTLrW2YadjCjfr6RWwSK76WxKGZ7uiT5QKDea0g7ggZvf1H6x1Sw/ZmCCDMsPXjowWqJaC7bZqk2Vjq/An09d6xMOXzL94L6Oy5FbkfBrkcUW48edtVACO44RPcN5DCA4tuQBfRG3luGUU5jWALcPMMoMYkrkSUIxPJTntxKXwrdyXg6gaXbTBQvmoozI3UVwwO0TKHnWzhcyFkLyzG9oopGZf2YtYu8/qUna7fZIj/3EUkv0supBywK176NKlXXZDv1mFhB+U/c1/YpbW34MxLc+vBTYwOj0/VCNjrSWNb3jsEHO1+18vyGrcbg1oT8PXDAd+qXR8+lXWeU2ZNIEDzYqhM1hSWofvvQ6Xv+4zv9ukohwMPZdmSVYvILl6+Ss011F4YVJTAhHpF3DQ3I25KJ2/dtDoTUY9Pp41U2V30sxTRq+hArmrwMyyiLIDi4xGpgbwae0ts82SKYdtPuc/APeP6lI8KFqEZiF+dO8ujzfYvqxPy6u3NBEThiyZrLwV+lRgiPZoMsgl6DI+zWdW13K6WKVyL3OtTuv2LrJGKcu/uU8+7gMnvgZY8ZgA70Qm5hgB"}],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":10,"cache_creation_input_tokens":30365,"cache_read_input_tokens":0,"cache_creation":{"ephemeral_5m_input_tokens":30365,"ephemeral_1h_input_tokens":0},"output_tokens":4,"service_tier":"standard"}},"parent_tool_use_id":null,"session_id":"3c350c77-fa41-4065-a987-bd07068b5127","uuid":"52b57e97-3728-4e35-97b4-278ea2adcecd"}"""
    l3 = r"""{"type":"assistant","message":{"id":"msg_01Lk3kuJspdcANLmH2jiXSpK","type":"message","role":"assistant","model":"claude-opus-4-1-20250805","content":[{"type":"text","text":"I'll begin the orchestrator continuation process with all mandatory checks and acknowledgments."}],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":10,"cache_creation_input_tokens":30365,"cache_read_input_tokens":0,"cache_creation":{"ephemeral_5m_input_tokens":30365,"ephemeral_1h_input_tokens":0},"output_tokens":4,"service_tier":"standard"}},"parent_tool_use_id":null,"session_id":"3c350c77-fa41-4065-a987-bd07068b5127","uuid":"471f19f6-d590-4cdb-b028-0a255b0c2562"}"""
    l4 = r"""{"type":"assistant","message":{"id":"msg_01Lk3kuJspdcANLmH2jiXSpK","type":"message","role":"assistant","model":"claude-opus-4-1-20250805","content":[{"type":"tool_use","id":"toolu_01Fnu8HArsDLJuRaNFWFMx2j","name":"Bash","input":{"command":"echo \"╔═══════════════════════════════════════════════════════════════════════════════╗\n║                        SOFTWARE FACTORY 2.0                                  ║\n║                   ORCHESTRATOR CONTINUATION                                   ║\n╚═══════════════════════════════════════════════════════════════════════════════╝\n\n🎯 AGENT IDENTITY: ORCHESTRATOR\n✓ Confirming identity: orchestrator\n\" && date '+%Y-%m-%d %H:%M:%S %Z'","description":"Start orchestrator continuation and show timestamp"}}],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":10,"cache_creation_input_tokens":30365,"cache_read_input_tokens":0,"cache_creation":{"ephemeral_5m_input_tokens":30365,"ephemeral_1h_input_tokens":0},"output_tokens":4,"service_tier":"standard"}},"parent_tool_use_id":null,"session_id":"3c350c77-fa41-4065-a987-bd07068b5127","uuid":"4d5a9313-6655-4c7d-8c07-c4365bb9dd80"}"""
    l5 = r"""{"type":"user","message":{"role":"user","content":[{"tool_use_id":"toolu_01Fnu8HArsDLJuRaNFWFMx2j","type":"tool_result","content":"╔═══════════════════════════════════════════════════════════════════════════════╗\n║                        SOFTWARE FACTORY 2.0                                  ║\n║                   ORCHESTRATOR CONTINUATION                                   ║\n╚═══════════════════════════════════════════════════════════════════════════════╝\n\n🎯 AGENT IDENTITY: ORCHESTRATOR\n✓ Confirming identity: orchestrator\n\n2025-09-23 16:22:45 UTC","is_error":false}]},"parent_tool_use_id":null,"session_id":"3c350c77-fa41-4065-a987-bd07068b5127","uuid":"0c5ba902-7fe2-436d-80fa-6b34821f9f64"}"""
    for line in (l1, l2, l3, l4, l5):
        print(line, flush=True)
    if args.emit_continuation:
        time.sleep(args.interval)
        print("CONTINUE-SOFTWARE-FACTORY=TRUE", flush=True)
else:
    # Simple mode: a short sequence of info and optional continuation
    print(json.dumps({
        "type": "system",
        "subtype": "init",
        "cwd": ".",
        "session_id": "mock-session",
        "tools": ["Task","Bash"],
        "model": "claude-opus-4-1-20250805",
        "permissionMode": "bypassPermissions",
        "output_style": "default"
    }), flush=True)

    print(json.dumps({
        "type": "assistant",
        "message": {"id": "msg_mock_1", "type": "message", "role": "assistant", "content": [{"type":"text","text":"I'll begin the orchestrator continuation process with all mandatory checks and acknowledgments."}]}
    }), flush=True)

    for i in range(3):
        time.sleep(args.interval)
        print(json.dumps({"type":"info","id": i, "message": f"tick {i}"}), flush=True)

    if args.emit_continuation:
        time.sleep(args.interval)
        print("CONTINUE-SOFTWARE-FACTORY=TRUE", flush=True)
