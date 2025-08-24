package repository

import (
	"database/sql"
	"github.com/Jersonmade/order-service-go/internal/model"
)

func SaveOrder(db *sql.DB, order *model.Order) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO orders (
			order_uid, 
			track_number, 
			entry, 
			locale, 
			internal_signature, 
			customer_id, 
			delivery_service, 
			shardkey, 
			sm_id, 
			date_created, 
			oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO deliveries (
			order_uid,
			name,
			phone,
			zip,
			city,
			address,
			region,
			email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO payments (
			order_uid,
			transaction,
			request_id,
			currency,
			provider,
			amount,
			payment_dt,
			bank,
			delivery_cost,
			goods_total,
			custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee)

	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (
				order_uid,
				chrt_id,
				track_number,
				price,
				rid,
				name,
				sale,
				size,
				total_price,
				nm_id,
				status,
				brand
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Status,
			item.Brand)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetOrderByUID(db *sql.DB, orderUID string) (*model.Order, error) {
	var order model.Order

	err := db.QueryRow(`SELECT order_uid, 
		track_number, 
		entry, 
		locale, 
		internal_signature, 
		customer_id, 
		delivery_service, 
		shardkey, 
		sm_id, 
		date_created, 
		oof_shard FROM orders WHERE order_uid = $1`, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard)

	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`SELECT name, phone, zip, city, address, region, email 
		FROM deliveries WHERE order_uid = $1`, orderUID).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email)

	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`SELECT transaction, 
		request_id, 
		currency, 
		provider, 
		amount, 
		payment_dt, 
		bank, 
		delivery_cost, 
		goods_total, 
		custom_fee FROM payments WHERE order_uid = $1`, orderUID).Scan(
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDT,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee)

	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status 
		FROM items WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}