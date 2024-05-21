go build -gcflags="all=-N -l" -o win-files-to-print-debug.exe
win-files-to-print-debug.exe install
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./win-files-to-print-debug.exe run
pause