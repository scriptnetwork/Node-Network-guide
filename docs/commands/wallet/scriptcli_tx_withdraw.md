## scriptcli tx withdraw

withdraw stake to a validator or guardian

### Synopsis

withdraw stake to a validator or guardian

```
scriptcli tx withdraw [flags]
```

### Examples

```
scriptcli tx withdraw --chain="scriptnet" --source=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --holder=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --purpose=0 --seq=8
```

### Options

```
      --chain string    Chain ID
      --fee string      Fee (default "1000000000000wei")
  -h, --help            help for withdraw
      --holder string   Holder of the stake
      --purpose uint8   Purpose of staking
      --seq uint        Sequence number of the transaction
      --source string   Source of the stake
      --wallet string   Wallet type (soft|nano) (default "soft")
```

### Options inherited from parent commands

```
      --config string   config path (default is /Users/<username>/.scriptcli) (default "/Users/<username>/.scriptcli")
```

### SEE ALSO

* [scriptcli tx](scriptcli_tx.md)	 - Manage transactions

###### Auto generated by spf13/cobra on 24-Apr-2019