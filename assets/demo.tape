Output assets/demo.gif

Set FontSize 32
Set Width 3000
Set Height 1500

Show
Type "cat assets/vops.yaml"
Enter
Sleep 5s

Hide
Type "clear"
Enter
Show 

Type "vops init --cluster cluster-1"
Enter
Sleep 5s

Type "vops unseal --cluster cluster-1"
Enter
Sleep 5s

Type "vops snapshot save --cluster cluster-1"
Enter
Sleep 5s

Type "vops config validate"
Enter
Sleep 5s