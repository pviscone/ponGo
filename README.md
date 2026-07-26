### Instructions

#### Build

```bash
go get example.com/ponGo
go build
./ponGo
```

or

```bash
go get example.com/ponGo
go run .
```

#### Play

Player1: Q UP/ A Down

Player2: O UP/ L Down

Key controls can be changed in `controls.go` Other game settings (ball speed,
etc.) can be changed in `main.go` in the `Game` struct initialization

```
  1                                                             0  
-------------------------------------------------------------------
                                                                   
                                                                   
                                                                   
                                                                   
                                                                   
                                                                   
%                                                                 %
%                                                               O %
%                                                                 %
                                                                   
                                                                   
                                                                   
                                                                   
                                                                   
                                                                   
-------------------------------------------------------------------
```
