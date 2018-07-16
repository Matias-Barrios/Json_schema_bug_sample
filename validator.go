package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/spf13/pflag"
	"encoding/json"
	"os"
	"io/ioutil"
	"strconv"
	"github.com/xeipuuv/gojsonschema"
)

const everything_is_ok = "Your file seems to be valid!"
var default_schema string = "./schemas/matias_schema.json"
var schemaLoader = gojsonschema.NewStringLoader(load_file(default_schema))

// This will be returned from the validate_spec function
	type validator_error struct{
		Path string
		Message string
	}
// ********************************************

func load_file(path string) string{
	b, err := ioutil.ReadFile(path) 
    if err != nil {
        fmt.Println(err)
		os.Exit(55)
    }
   return string(b) 
}

func validate_spec(spec string) []validator_error{
	var errors []validator_error
    
	api_spec := gojsonschema.NewStringLoader(spec)

    result, err := gojsonschema.Validate(schemaLoader, api_spec)
    if err != nil {
        panic(err.Error())
    }
	if ! result.Valid() {
      
        for _, err := range result.Errors() {
			n := validator_error{Path: err.Context().String(),  Message: err.Description()}
			errors = append(errors,n)
            
        }
	}
	
    return errors
}
func is_valid_JSON(s string) bool {
    var js map[string]interface{}
    return json.Unmarshal([]byte(s), &js) == nil

}
func main() {
		
			pflag.String("file", "", "Provide a valid path to a file to validate")
		    pflag.Parse()
			viper.BindPFlags(pflag.CommandLine)
			file_path := viper.GetString("file") 
		
			if (file_path != "") {
				
			
				file_contents := load_file(file_path)
				if (! is_valid_JSON(file_contents)) {
					fmt.Println("Your input is not a valid JSON!!");
					os.Exit(1)
				}
				
				errors := validate_spec(file_contents)
				for index,error := range errors {
					fmt.Println("----  Error " + strconv.Itoa(index) + " ----");
					fmt.Println("Path : " + error.Path);
					fmt.Println("Message : " + error.Message);
					fmt.Println("-----------------");
				}
							
		}
}

