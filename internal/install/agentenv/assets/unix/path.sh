# >>> traceknot-path >>>
if [ -d "{{BIN_DIR}}" ]; then
  case ":$PATH:" in
    *:"{{BIN_DIR}}":*) ;;
    *) export PATH="{{BIN_DIR}}:$PATH" ;;
  esac
fi
# <<< traceknot-path <<<
