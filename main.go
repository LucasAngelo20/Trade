package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	binanceSocketURL = "wss://stream.binance.com:9443/ws/btcusdt@ticker"
	bybitSocketURL   = "wss://stream.bybit.com/v5/public/linear"
)

type BBO struct {
	Market             string
	MarketHash         uint32
	Exchange           string
	ExchangeHash       uint32
	ExchangeType       string
	ExchangeTypeHash   uint32
	EntryEnabled       bool
	ExitEnabled        bool
	Quote              string
	QuoteHash          uint32
	Asset              string
	AssetHash          uint32
	Symbol             string
	SymbolHash         uint32
	Bid                float64
	Ask                float64
	BidMult            float64
	AskMult            float64
	BidSizeContracts   float64
	AskSizeContracts   float64
	BidSize            float64
	AskSize            float64
	TransactTime       int64
	UpdateTime         int64
	BuyEnabled         bool
	SellEnabled        bool
	BuyMultiplier      float64 // final multiplier for asks
	SellMultiplier     float64 // final multiplier for bids
	FeeMultiplier      float64 // fees from exchange
	OffsetMultiplier   float64 // offsets like averages, where an exchange is usually more expensive
	ContractMultiplier float64
	ContractSize       float64
	QtyDecimals        int
	PriceDecimals      int
	Timeout            bool
	WalletRef          *WalletData
	PositionRef        *PositionData
	// MidAssetRef        *EXCHANGE_MIDDLEWARE_ASSET
	// MidExchangeRef     *EXCHANGE_MIDDLEWARE_EXCHANGE
	// MidMarketRef       *EXCHANGE_MIDDLEWARE_MARKET
	// InstInfoRef        *InstrumentData
	MaxOrderSize       float64
	MinOrderSize       float64
	DefaultOrderData   *OrderData
	OgMsg              []byte
	UseProxy           bool
	OrderCount         *int
	Expiration         string
	ExpirationHash     uint32
	IsFuture           bool
	SocketDisconnected bool
	Leverage           float64
	Region             string
	RegionHash         uint32
}
type OrderData struct {
	Exchange               string
	ExchangeHash           uint32
	ExchangeType           string
	ExchangeTypeHash       uint32
	Symbol                 string `json:"symbol"`
	SymbolHash             uint32
	Asset                  string
	AssetHash              uint32
	Quote                  string
	QuoteHash              uint32
	Market                 string
	MarketHash             uint32
	AvgPrice               float64 `json:"avgPrice,omitempty,string"`
	AvgPriceMult           float64
	ClientOrderId          string  `json:"clientOrderId"`
	ClosePosition          bool    `json:"closePosition"`
	ExchangedQty           float64 `json:"exchangedQty,omitempty,string"`
	RemainingQty           float64 `json:"remainingQty,omitempty,string"`
	CumQty                 float64 `json:"cumQty,omitempty,string"`
	CumQuote               float64 `json:"cumQuote,omitempty,string"`
	ExecutedQty            float64 `json:"executedQty,omitempty,string"`
	OrderId                float64 `json:"orderId"`
	OrigQty                float64 `json:"origQty,omitempty,string"`
	OrigType               string  `json:"origType"`
	PositionSide           string  `json:"positionSide"`
	Price                  float64 `json:"price,omitempty,string"`
	PriceMult              float64
	PriceProtected         bool    `json:"priceProtected"`
	ReduceOnly             bool    `json:"reduceOnly"`
	Side                   string  `json:"side"`
	Status                 string  `json:"status"`
	StopPrice              float64 `json:"stopPrice,omitempty,string"`
	TimeInForce            string  `json:"timeInForce"`
	Type                   string  `json:"type"`
	UpdateTime             int64   `json:"updateTime"`
	WorkingType            string  `json:"workingType"`
	CancelReason           string
	ExecId                 string
	FakeOrder              bool
	Removing               bool
	TsClosed               int64
	Closed                 bool
	QuickOrder             bool
	ContractSizeMultiplier float64
	//DefaultOrderData map[uint32]*OrderData
	IsFuture bool
	BBORef   *BBO
}
type PositionData struct {
	ContractCode           string  `json:"contract_code"`
	AskNotional            float64 `json:"askNotional,omitempty,string"`
	BidNotional            float64 `json:"bidNotional,omitempty,string"`
	EntryPrice             float64 `json:"entryPrice,omitempty,string"`
	InitialMargin          float64 `json:"initialMargin,omitempty,string"`
	Isolated               bool    `json:"isolated"`
	IsolatedWallet         float64 `json:"isolatedWallet,omitempty,string"`
	Leverage               int     `json:"leverage,omitempty,string"`
	MaintMargin            float64 `json:"maintMargin,omitempty,string"`
	MaxNotional            float64 `json:"maxNotional,omitempty,string"`
	MaxQty                 float64 `json:"maxQty,omitempty,string"`
	Notional               float64 `json:"notional,omitempty,string"`
	NotionalValue          float64 `json:"notionalValue,omitempty,string"`
	OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,omitempty,string"`
	PositionAmt            float64 `json:"positionAmt,omitempty,string"`
	PositionInitialMargin  float64 `json:"positionInitialMargin,omitempty,string"`
	PositionSide           string  `json:"positionSide"`
	Symbol                 string  `json:"symbol"`
	SymbolHash             uint32
	Exchange               string
	ExchangeHash           uint32
	Quote                  string
	QuoteHash              uint32
	Asset                  string
	AssetHash              uint32
	UnrealizedProfit       float64 `json:"unrealizedProfit,omitempty,string"`
	UpdateTime             float64 `json:"updateTime"`
	Volume                 float64 `json:"volume"`
	SellVolume             float64
	BuyVolume              float64
	Direction              string `json:"direction"`
	Long                   bool
	NetPosition            float64
	NetPositionAbs         float64
	PositionMargin         float64 `json:"position_margin"`
	Profit                 float64 `json:"profit"`
	ProfitRate             float64 `json:"profit_rate"`
	ProfitUnreal           float64 `json:"profit_unreal"`
	SellPrice              float64
	Side                   string
	Short                  bool
	DualPosition           bool
	BBORef                 *BBO
}
type WalletData struct {
	Exchange               string
	ExchangeHash           uint32
	Symbol                 string
	SymbolHash             uint32
	Asset                  string `json:"asset"`
	AssetHash              uint32
	AvailableBalanceAux    float64 `json:"availableBalance,omitempty,string"`
	MaxWithdrawAmountAux   float64 `json:"maxWithdrawAmount,omitempty,string"`
	CrossUnPnl             float64 `json:"crossUnPnl,omitempty,string"`
	CrossWalletBalance     float64 `json:"crossWalletBalance,omitempty,string"`
	InitialMargin          float64 `json:"initialMargin,omitempty,string"`
	MaintMargin            float64 `json:"maintMargin,omitempty,string"`
	MarginAvailable        bool    `json:"marginAvailable"`
	MarginBalance          float64 `json:"marginBalance,omitempty,string"`
	TotalMarginBalance     float64 `json:"totalMarginBalance,omitempty,string"`
	OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,omitempty,string"`
	PositionInitialMargin  float64 `json:"positionInitialMargin,omitempty,string"`
	UnrealizedProfit       float64 `json:"unrealizedProfit,omitempty,string"`
	UpdateTime             float64 `json:"updateTime"`
	WalletBalance          float64 `json:"walletBalance,omitempty,string"`
	Free                   float64 `json:"free,omitempty,string"`
	Locked                 float64 `json:"locked,omitempty,string"`
	WithdrawAvailable      float64
	AccountId              float64 `json:"accountId"`
	AccountType            float64 `json:"accountType"`
	Currency               string  `json:"currency"`
	Balance                float64 `json:"balance,omitempty,string"`
	Available              float64 `json:"available,omitempty,string"`
	ChangeType             string  `json:"changeType"`
	ChangeTime             float64 `json:"changeTime"`
	SeqNum                 string  `json:"seqNum"`
	Equity                 float64
	MarginUsed             float64
	RiskExposure           float64
}

var BBOMap = make(map[string]*BBO)

type BBOL struct {
	Asset  string
	Quote  string
	Market string
}

var BBOLs = []BBOL{
	{"BTC", "USDT", "usdm"},
	// Add more symbols if needed
}

// connectToBinanceTicker connects to the Binance WebSocket ticker.
func connectToBinanceTicker() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(binanceSocketURL, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func connectToBybitTicker() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(bybitSocketURL, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// readTickerData reads ticker data from the WebSocket connection.
func readTickerData(conn *websocket.Conn, exchange string) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err)
			return
		}
		if exchange == "bybit" {

			if len(msg) > 120 {

				SerializeOrderbookDataLoop(msg)

			}

		}
	}
}

// SendTickerSubMsg sends the Bybit ticker subscription message.
func SendTickerSubMsg(connection *websocket.Conn) {
	subStr := GetTickerSubString(false)
	err := connection.WriteMessage(websocket.TextMessage, []byte(subStr))
	if err != nil {
		log.Printf("Error sending Bybit ticker sub message: %v\n", err)
		return
	}
}

// GetTickerSubString constructs the subscription message for Bybit ticker.
func GetTickerSubString(proxy bool) string {
	str := `{"req_id": "K2", "op": "subscribe", "args": [`
	for _, bbo := range BBOLs {
		s := strings.ToUpper(bbo.Asset + bbo.Quote)
		if bbo.Market == "usdm" {
			s = strings.ToUpper(bbo.Asset + ConvertUSDCToPERP(bbo.Quote))
		}
		str += "\"orderbook.1." + s + "\"" + ", "
	}
	str = str[:len(str)-2] + `],"id":1}`

	return str
}

func SerializeOrderbookDataLoop(message []byte) {
	// fmt.Printf("Price: %s \n", string(message))

	var response map[string]interface{}

	err := json.Unmarshal(message, &response)

	if err != nil {
		fmt.Printf("Bybit ticker socket serialize error: %s", err)

		return
	}

	data := response["data"].(map[string]interface{})

	// if data["success"] == false {
	// 	utils.AsyncPrint("red", "Bybit ticker socket serialize error: ")

	// 	return
	// }

	symbol := data["s"].(string)
	// bbo := BBOMap[symbol]
	bbo := *&BBO{Symbol: symbol}

	if len(data["b"].([]interface{})) != 0 {
		bidData := data["b"].([]interface{})
		for _, value := range bidData {
			bid := value.([]interface{})
			if bid[1].(string) != "0" {
				bbo.Bid = StringToFloat(bid[0].(string))
				bbo.BidSize = StringToFloat(bid[1].(string))
				break
			}
		}
	}

	bbo.BidMult = bbo.Bid * bbo.SellMultiplier

	if len(data["a"].([]interface{})) != 0 {
		askData := data["a"].([]interface{})
		for _, value := range askData {
			ask := value.([]interface{})
			if ask[1].(string) != "0" {
				bbo.Ask = StringToFloat(ask[0].(string))
				bbo.AskSize = StringToFloat(ask[1].(string))
				break
			}
		}

	}

	bbo.AskMult = bbo.Ask * bbo.BuyMultiplier

}

func main() {
	// Connect to Binance ticker
	binanceConn, err := connectToBinanceTicker()
	if err != nil {
		log.Fatal("Error connecting to Binance ticker:", err)
	}
	defer binanceConn.Close()

	// Connect to Bybit ticker
	bybitConn, err := connectToBybitTicker()
	if err != nil {
		log.Fatal("Error connecting to Bybit ticker:", err)
	}
	defer bybitConn.Close()

	// Send Bybit ticker subscription message
	SendTickerSubMsg(bybitConn)

	go readTickerData(bybitConn, "bybit")
	go readTickerData(binanceConn, "binance")

	// Keep the main goroutine alive
	for {
		time.Sleep(1 * time.Second)
	}
}

func ConvertUSDCToPERP(quote string) string {
	// Implement conversion logic if needed
	return quote
}
func StringToFloat(value string) float64 {
	newStr, _ := strconv.ParseFloat(value, 64)
	return newStr
}
