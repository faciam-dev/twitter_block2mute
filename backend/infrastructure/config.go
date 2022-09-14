package infrastructure

type Config struct {
    DB struct {
        Production struct {
            Host string
            Username string
            Password string
            DBName string
        }
        Test struct {
            Host string
            Username string
            Password string
            DBName string
        }
    }
    Routing struct {
        Port string
    }
    Session struct {
        Name string
        Secret string
    }
    Twitter struct {
        Production struct {
            ConsumerKey string
            ConsumerSecret string
            CallbackUrl string    
        }
        Test struct {
            ConsumerKey string
            ConsumerSecret string
            CallbackUrl string
        }
    }
}

func NewConfig() *Config {

    c := new(Config)

    // TODO: .envファイルまたはjsonファイルから読み込む
    c.DB.Production.Host = "localhost"
    c.DB.Production.Username = "username"
    c.DB.Production.Password = "password"
    c.DB.Production.DBName = "db_name"

    c.DB.Test.Host = "localhost"
    c.DB.Test.Username = "username"
    c.DB.Test.Password = "password"
    c.DB.Test.DBName = "db_name_test"
    
    c.Routing.Port = ":8080"

    c.Session.Name = "mysession"
    c.Session.Secret = "secret"

    c.Twitter.Production.CallbackUrl = ""
    c.Twitter.Production.ConsumerKey = ""
    c.Twitter.Production.ConsumerSecret = ""

    c.Twitter.Test.CallbackUrl = ""
    c.Twitter.Test.ConsumerKey = ""
    c.Twitter.Test.ConsumerSecret = ""

    return c
}