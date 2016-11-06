### wunderlist client


## example

```
func main() {
	cli, err := wlcli.NewClient("CLIENT-ID", "ACCESS-TOKEN")
	if err != nil {
		pp.Print(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, _ := cli.GetUser(ctx)

	pp.Println(user)

}
```