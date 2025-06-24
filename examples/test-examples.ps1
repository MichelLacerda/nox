Get-ChildItem -Path .\ -Recurse -Filter *.nox | ForEach-Object {
    Write-Host "Executing script:" $_.Name
    nox $_.FullName
}