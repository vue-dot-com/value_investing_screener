import yfinance as yf
import sys
import json

symbol = sys.argv[1]
action = sys.argv[2]

ticker = yf.Ticker(symbol)


def check_action_allowed(action: str) -> bool:
    attributes = [
        attr
        for attr in dir(ticker)
        if not attr.startswith("__") and not attr.startswith("_")
    ]

    # If the action is not a top-level attribute, parse the nested path
    if "[" in action:
        action_split = action.split("[")[0]  # Get the base attribute
        if action_split not in attributes:
            error_message = json.dumps(
                {
                    "error": f"Provided action '{action}' is not allowed. Allowed actions are: {attributes}"
                }
            )
            print(error_message, file=sys.stderr)  # Print error to stderr
            return False
    elif action not in attributes:
        error_message = json.dumps(
            {
                "error": f"Provided action '{action}' is not allowed. Allowed actions are: {attributes}"
            }
        )
        print(error_message, file=sys.stderr)  # Print error to stderr
        return False

    return True


def resolve_nested_attr(obj: classmethod, path: str) -> str:
    """Recursively resolve nested attributes and dictionary keys"""
    parts = path.split("[")
    base_attr = parts[0]
    nested_keys = [
        p.rstrip("]") for p in parts[1:]
    ]  # Clean the keys (remove trailing ']')

    # Get the base attribute
    result = getattr(obj, base_attr, None)

    # Resolve nested keys
    for key in nested_keys:
        if key in result.keys():
            result = result[key]
        else:
            raise AttributeError(f"Key '{key}' not found in {base_attr} attribute.")

    return result


if check_action_allowed(action):
    try:
        # If action contains a nested key (like 'fast_info['lastPrice']'), resolve it
        if "[" in action:
            result = resolve_nested_attr(ticker, action)
        else:
            result = getattr(ticker, action)

        # Debugging statement
        # print(f"Resolved result for action '{action}': {result}")

        # Ensure result is a string for JSON serialization
        if result is None:
            result = "None"

        print(json.dumps({"symbol": symbol, "action": action, "result": str(result)}))
        sys.stdout.flush()  # Ensure stdout is flushed
    except Exception as e:
        error_message = json.dumps(
            {"error": f"Error occurred for action '{action}': {str(e)}"}
        )
        print(error_message, file=sys.stderr)  # Print error to stderr
        sys.stderr.flush()  # Ensure stderr is flushed
else:
    print(json.dumps({"error": "Invalid action provided."}), file=sys.stderr)
    sys.stderr.flush()
