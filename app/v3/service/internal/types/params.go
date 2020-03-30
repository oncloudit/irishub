package types

import (
	"fmt"
	"strconv"
	"time"

	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

var _ params.ParamSet = (*Params)(nil)

// default paramSpace for service keeper
const (
	DefaultParamSpace = "service"
)

//Parameter store key
var (
	// params store for service params
	KeyMaxRequestTimeout    = []byte("MaxRequestTimeout")
	KeyMinDepositMultiple   = []byte("MinDepositMultiple")
	KeyMinDeposit           = []byte("MinDeposit")
	KeyServiceFeeTax        = []byte("ServiceFeeTax")
	KeySlashFraction        = []byte("SlashFraction")
	KeyComplaintRetrospect  = []byte("ComplaintRetrospect")
	KeyArbitrationTimeLimit = []byte("ArbitrationTimeLimit")
	KeyTxSizeLimit          = []byte("TxSizeLimit")
)

// ParamTable for service module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// service params
type Params struct {
	MaxRequestTimeout    int64         `json:"max_request_timeout"`
	MinDepositMultiple   int64         `json:"min_deposit_multiple"`
	MinDeposit           sdk.Coins     `json:"min_deposit"`
	ServiceFeeTax        sdk.Dec       `json:"service_fee_tax"`
	SlashFraction        sdk.Dec       `json:"slash_fraction"`
	ComplaintRetrospect  time.Duration `json:"complaint_retrospect"`
	ArbitrationTimeLimit time.Duration `json:"arbitration_time_limit"`
	TxSizeLimit          uint64        `json:"tx_size_limit"`
}

func (p Params) String() string {
	return fmt.Sprintf(`Service Params:
  service/MaxRequestTimeout:     %d
  service/MinDepositMultiple:    %d
  service/MinDeposit:            %s
  service/ServiceFeeTax:         %s
  service/SlashFraction:         %s
  service/ComplaintRetrospect:   %s
  service/ArbitrationTimeLimit:  %s
  service/TxSizeLimit:           %d`,
		p.MaxRequestTimeout, p.MinDepositMultiple, p.MinDeposit.String(), p.ServiceFeeTax.String(), p.SlashFraction.String(),
		p.ComplaintRetrospect, p.ArbitrationTimeLimit, p.TxSizeLimit)
}

// Implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyMaxRequestTimeout, Value: &p.MaxRequestTimeout},
		{Key: KeyMinDepositMultiple, Value: &p.MinDepositMultiple},
		{Key: KeyMinDeposit, Value: &p.MinDeposit},
		{Key: KeyServiceFeeTax, Value: &p.ServiceFeeTax},
		{Key: KeySlashFraction, Value: &p.SlashFraction},
		{Key: KeyComplaintRetrospect, Value: &p.ComplaintRetrospect},
		{Key: KeyArbitrationTimeLimit, Value: &p.ArbitrationTimeLimit},
		{Key: KeyTxSizeLimit, Value: &p.TxSizeLimit},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyMaxRequestTimeout):
		maxRequestTimeout, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMaxRequestTimeout(maxRequestTimeout); err != nil {
			return nil, err
		}
		return maxRequestTimeout, nil
	case string(KeyMinDepositMultiple):
		minDepositMultiple, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMinDepositMultiple(minDepositMultiple); err != nil {
			return nil, err
		}
		return minDepositMultiple, nil
	case string(KeyMinDeposit):
		minDeposit, err := sdk.ParseCoins(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMinDeposit(minDeposit); err != nil {
			return nil, err
		}
		return minDeposit, nil
	case string(KeyServiceFeeTax):
		serviceFeeTax, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateServiceFeeTax(serviceFeeTax); err != nil {
			return nil, err
		}
		return serviceFeeTax, nil
	case string(KeySlashFraction):
		slashFraction, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateSlashFraction(slashFraction); err != nil {
			return nil, err
		}
		return slashFraction, nil
	case string(KeyComplaintRetrospect):
		complaintRetrospect, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateComplaintRetrospect(complaintRetrospect); err != nil {
			return nil, err
		}
		return complaintRetrospect, nil
	case string(KeyArbitrationTimeLimit):
		arbitrationTimeLimit, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateArbitrationTimeLimit(arbitrationTimeLimit); err != nil {
			return nil, err
		}
		return arbitrationTimeLimit, nil
	case string(KeyTxSizeLimit):
		txSizeLimit, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateTxSizeLimit(txSizeLimit); err != nil {
			return nil, err
		}
		return txSizeLimit, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeyMaxRequestTimeout):
		err := cdc.UnmarshalJSON(bytes, &p.MaxRequestTimeout)
		return strconv.FormatInt(p.MaxRequestTimeout, 10), err
	case string(KeyMinDepositMultiple):
		err := cdc.UnmarshalJSON(bytes, &p.MinDepositMultiple)
		return strconv.FormatInt(p.MinDepositMultiple, 10), err
	case string(KeyMinDeposit):
		err := cdc.UnmarshalJSON(bytes, &p.MinDeposit)
		return p.MinDeposit.String(), err
	case string(KeyServiceFeeTax):
		err := cdc.UnmarshalJSON(bytes, &p.ServiceFeeTax)
		return p.ServiceFeeTax.String(), err
	case string(KeySlashFraction):
		err := cdc.UnmarshalJSON(bytes, &p.SlashFraction)
		return p.SlashFraction.String(), err
	case string(KeyComplaintRetrospect):
		err := cdc.UnmarshalJSON(bytes, &p.ComplaintRetrospect)
		return p.ComplaintRetrospect.String(), err
	case string(KeyArbitrationTimeLimit):
		err := cdc.UnmarshalJSON(bytes, &p.ArbitrationTimeLimit)
		return p.ArbitrationTimeLimit.String(), err
	case string(KeyTxSizeLimit):
		err := cdc.UnmarshalJSON(bytes, &p.TxSizeLimit)
		return strconv.FormatUint(p.TxSizeLimit, 10), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

func (p *Params) ReadOnly() bool {
	return false
}

// default service module params
func DefaultParams() Params {
	return Params{
		MaxRequestTimeout:    100,
		MinDepositMultiple:   1000,
		MinDeposit:           sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewIntWithDecimal(10, 18))), // 10000iris
		ServiceFeeTax:        sdk.NewDecWithPrec(1, 2),                                               // 1%
		SlashFraction:        sdk.NewDecWithPrec(1, 3),                                               // 0.1%
		ComplaintRetrospect:  time.Duration(15 * sdk.Day),                                            // 15 days
		ArbitrationTimeLimit: time.Duration(5 * sdk.Day),                                             // 5 days
		TxSizeLimit:          4000,
	}
}

// default service module params for test
func DefaultParamsForTest() Params {
	return Params{
		MaxRequestTimeout:    10,
		MinDepositMultiple:   10,
		MinDeposit:           sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewIntWithDecimal(10, 18))), // 10iris
		ServiceFeeTax:        sdk.NewDecWithPrec(1, 2),                                               // 1%
		SlashFraction:        sdk.NewDecWithPrec(1, 3),                                               // 0.1%
		ComplaintRetrospect:  20 * time.Second,                                                       // 20 seconds
		ArbitrationTimeLimit: 20 * time.Second,                                                       // 20 seconds
		TxSizeLimit:          4000,
	}
}

func validateParams(p Params) error {
	if err := validateMaxRequestTimeout(p.MaxRequestTimeout); err != nil {
		return err
	}
	if err := validateMinDepositMultiple(p.MinDepositMultiple); err != nil {
		return err
	}
	if err := validateMinDeposit(p.MinDeposit); err != nil {
		return err
	}
	if err := validateSlashFraction(p.SlashFraction); err != nil {
		return err
	}
	if err := validateServiceFeeTax(p.ServiceFeeTax); err != nil {
		return err
	}
	if err := validateComplaintRetrospect(p.ComplaintRetrospect); err != nil {
		return err
	}
	if err := validateArbitrationTimeLimit(p.ArbitrationTimeLimit); err != nil {
		return err
	}
	if err := validateTxSizeLimit(p.TxSizeLimit); err != nil {
		return err
	}
	return nil
}

func validateMaxRequestTimeout(v int64) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if v < 20 {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxRequestTimeout, fmt.Sprintf("Invalid MaxRequestTimeout [%d] should be greater than or equal to 20", v))
		}
	} else if v < 5 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxRequestTimeout, fmt.Sprintf("Invalid MaxRequestTimeout [%d] should be greater than or equal to 5", v))
	}
	return nil
}

func validateMinDepositMultiple(v int64) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if v < 500 || v > 5000 {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositMultiple, fmt.Sprintf("Invalid MinDepositMultiple [%d] should be between [500, 5000]", v))
		}
	} else if v < 10 || v > 5000 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositMultiple, fmt.Sprintf("Invalid MinDepositMultiple [%d] should be between [10, 5000]", v))
	}
	return nil
}

func validateMinDeposit(v sdk.Coins) sdk.Error {
	max := sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewIntWithDecimal(100000, 18)))
	var min sdk.Coins

	if sdk.NetworkType == sdk.Mainnet {
		min = sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewIntWithDecimal(5000, 18)))
	} else {
		min = sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewIntWithDecimal(10, 18)))
	}

	if v.IsAllLT(min) || v.IsAllGT(max) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDeposit, fmt.Sprintf("Invalid MinDeposit [%s] should be between [%s, %s]", v, min, max))
	}

	return nil
}

func validateSlashFraction(v sdk.Dec) sdk.Error {
	if v.LTE(sdk.ZeroDec()) || v.GT(sdk.NewDecWithPrec(1, 2)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashFraction, fmt.Sprintf("Invalid SlashFraction [%s] should be between (0, 0.01]", v.String()))
	}
	return nil
}

func validateServiceFeeTax(v sdk.Dec) sdk.Error {
	if v.LTE(sdk.ZeroDec()) || v.GT(sdk.NewDecWithPrec(2, 1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidServiceFeeTax, fmt.Sprintf("Invalid ServiceFeeTax [%s] should be between (0, 0.2]", v.String()))
	}
	return nil
}

func validateComplaintRetrospect(v time.Duration) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if v < 15*sdk.Day || v > 30*sdk.Day {
			return sdk.NewError(params.DefaultCodespace, params.CodeComplaintRetrospect, fmt.Sprintf("Invalid ComplaintRetrospect [%s] should be between [15days, 30days]", v.String()))
		}
	} else if v < 20*time.Second {
		return sdk.NewError(params.DefaultCodespace, params.CodeComplaintRetrospect, fmt.Sprintf("Invalid ComplaintRetrospect [%s] should be between [20seconds, )", v.String()))
	}
	return nil
}

func validateArbitrationTimeLimit(v time.Duration) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if v < 5*sdk.Day || v > 10*sdk.Day {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidArbitrationTimeLimit, fmt.Sprintf("Invalid ArbitrationTimeLimit [%s] should be between [5days, 10days]", v.String()))
		}
	} else if v < 20*time.Second {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidArbitrationTimeLimit, fmt.Sprintf("Invalid ArbitrationTimeLimit [%s] should be between [20seconds, )", v.String()))
	}
	return nil
}

func validateTxSizeLimit(v uint64) sdk.Error {
	if v < uint64(2000) || v > uint64(6000) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidServiceTxSizeLimit, fmt.Sprintf("Invalid ServiceTxSizeLimit [%d] should be between [2000, 6000]", v))
	}
	return nil
}
