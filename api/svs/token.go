package svs

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// TokenChgByPeriod gorm columns are named with underscores, json columns are named with camelCase

var sqlQueryPriceChange = `
WITH time_periods AS (
  SELECT 
    '5m' AS period,
    DATE_SUB(NOW() %s, INTERVAL 5 MINUTE) AS start_time,
    NOW() AS end_time
  UNION ALL
  SELECT 
    '1h' AS period,
    DATE_SUB(NOW() %s, INTERVAL 1 HOUR) AS start_time,
    NOW() AS end_time
  UNION ALL
  SELECT 
    '4h' AS period,
    DATE_SUB(NOW() %s, INTERVAL 4 HOUR) AS start_time,
    NOW() AS end_time
  UNION ALL
  SELECT 
    '1d' AS period,
    DATE_SUB(NOW() %s, INTERVAL 1 DAY) AS start_time,
    NOW() AS end_time
),
ranked_prices AS (
  SELECT 
    stf.token0,
    stf.price,
    stf.tx_time,  -- 添加 tx_time 以便调试或验证
    tp.period,
    ROW_NUMBER() OVER (PARTITION BY stf.token0, tp.period ORDER BY stf.tx_time ASC) AS rn_asc,
    ROW_NUMBER() OVER (PARTITION BY stf.token0, tp.period ORDER BY stf.tx_time DESC) AS rn_desc
  FROM 
    sol_token_flow stf
  JOIN 
    time_periods tp
  ON 
    stf.tx_time >= tp.start_time 
    AND stf.tx_time <= tp.end_time
  WHERE 
    stf.token0 = '%s'
),
price_changes AS (
  SELECT 
    period,
    MAX(CASE WHEN rn_asc = 1 THEN price END) AS start_price,
    MAX(CASE WHEN rn_desc = 1 THEN price END) AS end_price,
    MAX(CASE WHEN rn_asc = 1 THEN tx_time END) AS start_time,  -- 添加对应的开始时间
    MAX(CASE WHEN rn_desc = 1 THEN tx_time END) AS end_time    -- 添加对应的结束时间
  FROM 
    ranked_prices
  GROUP BY 
    period
)
SELECT 
  period,
  start_price,
  end_price,
  start_time,
  end_time
FROM 
  price_changes
ORDER BY 
  period;`

var sqlQueryR = `WITH time_periods AS (
  SELECT 
    '5m' AS period,
    DATE_SUB(NOW() %s, INTERVAL 5 MINUTE) AS start_time,
    NOW() AS end_time
  UNION ALL
  SELECT 
    '1h' AS period,
    DATE_SUB(NOW() %s, INTERVAL 1 HOUR) AS start_time,
    NOW() AS end_time
  UNION ALL
  SELECT 
    '4h' AS period,
    DATE_SUB(NOW() %s, INTERVAL 4 HOUR) AS start_time,
    NOW() AS end_time
  UNION ALL
  SELECT 
    '1d' AS period,
    DATE_SUB(NOW() %s, INTERVAL 1 DAY) AS start_time,
    NOW() AS end_time
),
ranked_prices AS (
  SELECT 
    stf.token0,
    stf.price,
    tp.period,
    ROW_NUMBER() OVER (PARTITION BY stf.token0, tp.period ORDER BY stf.tx_time ASC) AS rn_asc,
    ROW_NUMBER() OVER (PARTITION BY stf.token0, tp.period ORDER BY stf.tx_time DESC) AS rn_desc
  FROM 
    sol_token_flow stf
  JOIN 
    time_periods tp
  ON 
    stf.tx_time >= tp.start_time 
    AND stf.tx_time <= tp.end_time
  WHERE 
    stf.token0 = '%s'
),
price_changes AS (
  SELECT 
    period,
    MAX(CASE WHEN rn_asc = 1 THEN price END) AS start_price,
    MAX(CASE WHEN rn_desc = 1 THEN price END) AS end_price
  FROM 
    ranked_prices
  GROUP BY 
    period
),
trade_stats AS (
  SELECT 
    tp.period,
    COUNT(stf.tx) AS tradeCount,
    SUM(CASE WHEN stf.amount0_in > 0 THEN 1 ELSE 0 END) AS tradeBuyCount,
    SUM(CASE WHEN stf.amount0_out > 0 THEN 1 ELSE 0 END) AS tradeSellCount,
    SUM(CASE WHEN stf.amount0_in > 0 THEN CAST(stf.amount0_in AS DECIMAL(32,6)) / POW(10, stf.token0_decimals) * stf.price ELSE 0 END) AS tradeBuy,
    SUM(CASE WHEN stf.amount0_out > 0 THEN CAST(stf.amount0_out AS DECIMAL(32,6)) / POW(10, stf.token0_decimals) * stf.price ELSE 0 END) AS tradeSell,
    SUM(CAST(stf.amount0_in + stf.amount0_out AS DECIMAL(32,6)) / POW(10, stf.token0_decimals) * stf.price) AS tradeVol
  FROM 
    sol_token_flow stf
  JOIN 
    time_periods tp
  ON 
    stf.tx_time >= tp.start_time 
    AND stf.tx_time <= tp.end_time
  WHERE 
    stf.token0 = '%s'
  GROUP BY 
    tp.period
),
valid_periods AS (
  SELECT period
  FROM trade_stats ts
  WHERE ts.tradeCount > 0
)
SELECT 
  ts.period,
  pc.start_price,
  pc.end_price,
  ts.tradeCount,
  ts.tradeBuyCount,
  ts.tradeSellCount,
  ts.tradeBuy,
  ts.tradeSell,
  ts.tradeVol
FROM 
  trade_stats ts
LEFT JOIN 
  price_changes pc ON ts.period = pc.period
JOIN 
  valid_periods vp ON ts.period = vp.period;  
`
