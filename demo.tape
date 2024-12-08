Output demo.gif

Set Shell "bash"
Set FontSize 20
Set Width 1200
Set Height 600

Type "cd ~/documents/go-projects/sqlite-shell" Sleep 500ms
Enter Sleep 500ms

Type "go run . ~/chinook.db" Sleep 500ms
Enter Sleep 4s

Type 'SELECT BillingCountry, SUM(Total) FROM invoice' Sleep 500ms
Enter Sleep 2s

Type 's GROUP BY BillingCountry' Sleep 500ms
Enter Sleep 2s

Down 12 Sleep 2s


