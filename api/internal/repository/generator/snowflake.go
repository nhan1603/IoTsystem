package generator

import "github.com/nhan1603/IoTsystem/api/internal/pkg/snowflake"

var (
	DeviceTokenIDSNF snowflake.SnowflakeGenerator
	AlertIDSNF       snowflake.SnowflakeGenerator
	RequestIDSNF     snowflake.SnowflakeGenerator
	ResponseSNF      snowflake.SnowflakeGenerator
)

func InitSnowflakeGenerators() error {
	DeviceTokenIDSNF = snowflake.New()
	AlertIDSNF = snowflake.New()
	RequestIDSNF = snowflake.New()
	ResponseSNF = snowflake.New()
	return nil
}
