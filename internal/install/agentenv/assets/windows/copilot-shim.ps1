# >>> traceknot-copilot >>>
function copilot {
    & "{{TRACEKNOT_BIN}}" prep-copilot
    $real = Get-Command copilot -All -ErrorAction SilentlyContinue |
        Where-Object { $_.CommandType -ne 'Function' } |
        Select-Object -First 1
    if ($real) {
        & $real.Source @args
    }
}
# <<< traceknot-copilot <<<
