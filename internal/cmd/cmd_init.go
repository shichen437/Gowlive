package cmd

import (
	_ "github.com/shichen437/gowlive/internal/app/admin/logic"
	_ "github.com/shichen437/gowlive/internal/app/common/logic"
	_ "github.com/shichen437/gowlive/internal/app/media/logic"
	_ "github.com/shichen437/gowlive/internal/app/stream/logic"
	_ "github.com/shichen437/gowlive/internal/app/system/logic"

	_ "github.com/shichen437/gowlive/internal/pkg/lives/bilibili"
	_ "github.com/shichen437/gowlive/internal/pkg/lives/douyin"
	_ "github.com/shichen437/gowlive/internal/pkg/lives/yy"

	_ "github.com/shichen437/gowlive/internal/pkg/anchor/bilibili"
	_ "github.com/shichen437/gowlive/internal/pkg/anchor/douyin"
	_ "github.com/shichen437/gowlive/internal/pkg/anchor/yy"

	_ "github.com/shichen437/gowlive/internal/pkg/stream_parser/ffmpeg"

	_ "github.com/shichen437/gowlive/internal/pkg/message_push/ding_talk"
	_ "github.com/shichen437/gowlive/internal/pkg/message_push/email"
	_ "github.com/shichen437/gowlive/internal/pkg/message_push/gotify"
	_ "github.com/shichen437/gowlive/internal/pkg/message_push/lark"
	_ "github.com/shichen437/gowlive/internal/pkg/message_push/we_com"

	_ "github.com/shichen437/gowlive/internal/pkg/monitor"
)
