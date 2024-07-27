package main

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type OmsEmail struct {
	Id                int         `json:"id"`
	Email             string      `json:"email"`
	User_id           int         `json:"user_id"`
	Is_subscribe      int         `json:"is_subscribe"`
	Is_resubscribe    int         `json:"is_resubscribe"`
	Language          string      `json:"language"`
	Site              string      `json:"site"`
	Subscribe_channel string      `json:"subscribe_channel"`
	Platform          string      `json:"platform"`
	Created_at        string      `json:"created_at"`
	Updated_at        string      `json:"updated_at"`
	Deleted_at        interface{} `json:"deleted_at"`
}

// 查询该email是否已经存在
func findEmail(email string) (bool, error) {
	var exists bool
	err := readonlyDB.QueryRow("SELECT 1 FROM oms_emails WHERE email = ?", email).Scan(&exists)
	if err != nil {
		return false, err
	} else {
		if err == sql.ErrNoRows {
			return false, err
		}
	}
	return exists, err
}
func insertEmailRecord(record OmsEmail) (sql.Result, error) {
	query := `INSERT INTO oms_emails (
			email, user_id, is_subscribe, is_resubscribe, language, site,
			subscribe_channel, platform, created_at, updated_at, deleted_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := writeDB.Exec(query, record.Email, record.User_id, record.Is_subscribe,
		record.Is_resubscribe, record.Language, record.Site, record.Subscribe_channel,
		record.Platform, record.Created_at, record.Updated_at, record.Deleted_at)

	return res, err
}

func getOmsEmails(c *fiber.Ctx) error {
	var omsEmailsResponse []OmsEmail
	rows, err := readonlyDB.Query("SELECT * FROM oms_emails LIMIT 10")
	if err != nil {
		return createResponseError(c, 400, err.Error())
	}
	defer rows.Close()
	// 处理查询结果
	for rows.Next() {
		var omsEmail OmsEmail

		// 根据实际表结构替换
		err := rows.Scan(&omsEmail.Id, &omsEmail.Email, &omsEmail.User_id, &omsEmail.Is_subscribe, &omsEmail.Is_resubscribe, &omsEmail.Language, &omsEmail.Site, &omsEmail.Subscribe_channel, &omsEmail.Platform, &omsEmail.Created_at, &omsEmail.Updated_at, &omsEmail.Deleted_at) // 根据列数和类型调整
		// 如果是空，转换空字符串
		if omsEmail.Deleted_at == nil {
			omsEmail.Deleted_at = ""
		}
		omsEmailsResponse = append(omsEmailsResponse, omsEmail)
		if err != nil {
			createResponseError(c, 400, err.Error())
			break
		}
	}

	// 检查遍历过程中的错误
	if err = rows.Err(); err != nil {
		return createResponseError(c, 400, err.Error())
	}

	return createResponseSuccess(c, omsEmailsResponse)
}

func getOmsEmailDetail(c *fiber.Ctx) error {
	var omsEmail OmsEmail
	email := c.Query("email")
	rows := readonlyDB.QueryRow("SELECT * FROM oms_emails WHERE email = ?", email)
	err := rows.Scan(&omsEmail.Id, &omsEmail.Email, &omsEmail.User_id, &omsEmail.Is_subscribe, &omsEmail.Is_resubscribe, &omsEmail.Language, &omsEmail.Site, &omsEmail.Subscribe_channel, &omsEmail.Platform, &omsEmail.Created_at, &omsEmail.Updated_at, &omsEmail.Deleted_at) // 根据列数和类型调整
	if err != nil {
		return createResponseError(c, 400, err.Error())
	}
	return createResponseSuccess(c, omsEmail)
}

func updataOmsEmail(c *fiber.Ctx) error {
	email := c.Query("email")

	log.Info(email)
	if email == "" {
		return createResponseError(c, 400, "请输入email!")
	}
	var isSubscribe int
	rows := readonlyDB.QueryRow("SELECT is_subscribe FROM oms_emails WHERE email = ?", email)
	err := rows.Scan(&isSubscribe)
	log.Info(isSubscribe)
	if err != nil {
		if err == sql.ErrNoRows {
			return createResponseError(c, 400, "用户不存在")
		}
		return createResponseError(c, 400, err.Error())
	}
	newIsSubscribe := 1
	if isSubscribe == 1 {
		newIsSubscribe = 0
	}
	start := time.Now()
	result, err := writeDB.Exec("UPDATE oms_emails SET is_subscribe = ? WHERE email = ?", newIsSubscribe, email)
	if err != nil {
		return createResponseError(c, 200, err.Error())
	}

	// 获取受影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return createResponseError(c, 200, err.Error())
	}

	log.Info("更新了", rowsAffected)
	log.Infof("更新操作耗时: %v\n", time.Since(start))
	// writeDB.Exec("UPDATE oms_emails SET is_subscribe = ? WHERE email = ?", newIsSubscribe, email)
	// errorRes, _ := json.Marshal(Response{
	// 	Message: "Updata Success!",
	// 	Code:    fiber.StatusOK,
	// 	Data:    nil,
	// })

	return createResponseSuccess(c, nil, "Updata Success!")
}

/*
	 @desc 设置邮箱信息
		Id                int         `json:"id"`
		Email             string      `json:"email"`
		User_id           int         `json:"user_id"`
		Is_subscribe      int         `json:"is_subscribe"`
		Is_resubscribe    int         `json:"is_resubscribe"`
		Language          string      `json:"language"`
		Site              string      `json:"site"`
		Subscribe_channel string      `json:"subscribe_channel"`
		Platform          string      `json:"platform"`
		Created_at        string      `json:"created_at"`
		Updated_at        string      `json:"updated_at"`
		Deleted_at        interface{} `json:"deleted_at"`
*/
func setOmsEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if !(isString(email)) {
		return createResponseError(c, 400, "请输入email")
	}
	if !(isEmail(email)) {
		return createResponseError(c, 400, "email 无效")
	}
	// 查询该email是否已经存在
	exists, _ := findEmail(email)
	if exists {
		return createResponseError(c, 400, "用户已存在")
	}
	if !(isEmail(email)) {
		return createResponseError(c, 400, "email 无效")
	}
	userId, userIdRrr := strconv.Atoi(c.Query("user_id"))
	if userIdRrr != nil || userId == 0 {
		return createResponseError(c, 400, "请输入正确的user_id")
	}
	// if !(isInt(userId)) {
	// 	return createResponseError(c, 400, "user_id 无效")
	// }
	// is_resubscribe

	isSubscribe, isSubscribeError := strconv.Atoi(c.Query("is_subscribe"))
	if isSubscribeError != nil {
		isSubscribe = 0
	}

	isResubscribe, isResubscribeError := strconv.Atoi(c.Query("is_resubscribe"))
	if isResubscribeError != nil {
		isResubscribe = 0
	}
	language := c.Query("language")
	if language == "" {
		language = "en"
	}
	site := c.Query("site")
	if site == "" {
		site = "us"
	}
	platform := c.Query("platform")
	if platform == "" {
		platform = "web"
	}
	subscribeChannel := c.Query("subscribe_channel")
	now := time.Now().Format("2006-01-02 15:04:05")
	omsEmail := OmsEmail{
		Email:             email,
		User_id:           userId,
		Is_subscribe:      isSubscribe,
		Is_resubscribe:    isResubscribe,
		Language:          language,
		Site:              site,
		Subscribe_channel: subscribeChannel,
		Platform:          platform,
		Created_at:        now,
		Updated_at:        now,
		Deleted_at:        nil,
	}
	res, err := insertEmailRecord(omsEmail)
	if err != nil {
		createResponseError(c, 400, err.Error())
	}
	log.Info(res)
	return createResponseSuccess(c, email, "Add Success")
}

func deleteOmsEmailByEmail(c *fiber.Ctx) error {
	// 如果是需要权限操作，先鉴权再删除

	email := c.Query("email")
	if !(isString(email)) {
		return createResponseError(c, 400, "请输入email")
	}
	if !(isEmail(email)) {
		return createResponseError(c, 400, "email 无效")
	}

	// 查询该email是否已经存在
	exists, _ := findEmail(email)
	if !exists {
		return createResponseError(c, 400, "用户不存在")
	}
	query := `DELETE FROM oms_emails WHERE email = ?`

	_, err := writeDB.Exec(query, email)
	if err != nil {
		return createResponseError(c, 400, "删除失败")
	}

	// result, err := writeDB.Exec(query, email)
	// // 获取受影响的行数
	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	return createResponseError(c, 400, "删除失败")
	// }

	return createResponseSuccess(c, nil, "删除成功")
}
