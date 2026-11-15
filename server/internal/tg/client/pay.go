package client

import "context"

type Price struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

func (c *Client) Invoice(__ctx context.Context, __title, __desc, __payload string, __amount, __period int) (string, error) {
	var _url string
	_err := c.Call(__ctx, "createInvoiceLink", map[string]any{
		"title":               __title,
		"description":         __desc,
		"payload":             __payload,
		"currency":            "XTR",
		"prices":              []Price{{Label: __title, Amount: __amount}},
		"subscription_period": __period,
	}, &_url)
	return _url, _err
}

func (c *Client) PreOk(__ctx context.Context, __id string, __err string) error {
	_args := map[string]any{"pre_checkout_query_id": __id, "ok": __err == ""}
	if __err != "" {
		_args["error_message"] = __err
	}
	return c.Call(__ctx, "answerPreCheckoutQuery", _args, nil)
}

func (c *Client) Cancel(__ctx context.Context, __user int64, __charge string) error {
	return c.Call(__ctx, "editUserStarSubscription", map[string]any{
		"user_id": __user, "telegram_payment_charge_id": __charge, "is_canceled": true,
	}, nil)
}
