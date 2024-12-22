package accounts

type Accounts struct {
	tableName struct{} `table:"accounts"`
	ID        int      `field:"id"`
	Username  string   `field:"username"`
}
