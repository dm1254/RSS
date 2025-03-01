package config
import(
	"encoding/json"
	"os"
	"io"
)

type Config struct{
	DBURL string `json:"db_url"`
	Username string `json:"current_user_name"`
	
}

func(c *Config) SetUser(username string) error{
	home_dir, err := os.UserHomeDir()
	if err != nil{
		return err
	}
	c.Username = username 
	jsonData, err := json.Marshal(c)
	if err != nil{
		return err
	}
	if err = os.WriteFile(home_dir + "/.gatorconfig.json", jsonData,0644); err != nil{
		return err
	}
	return nil
}	

func Read() (Config,error){
	home_dir, err := os.UserHomeDir()
	if err != nil{
		return Config{},err
	}
	file,err := os.Open(home_dir + "/.gatorconfig.json")
	if err != nil{
		return Config{},err
	}
	defer file.Close()

	data,err := io.ReadAll(file)
	if err != nil{
		return Config{},err
	}
	var config Config
	if err = json.Unmarshal(data, &config); err != nil{
		return Config{},err
	}
	return config,nil

}
