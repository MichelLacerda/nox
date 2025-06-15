```powershell
Get-ChildItem -Recurse -Filter *.go | ForEach-Object {
    Write-Host "`n==== $($_.FullName) ====" -ForegroundColor Cyan
    Get-Content $_
}
```