function genProto {
   DOMAIN=$1
    PROTO_PATH=./${DOMAIN}/api
    GO_OUT_PATH=./${DOMAIN}/api/gen/v1
    protoc -I=$PROTO_PATH --go_out=plugins=grpc,paths=source_relative:$GO_OUT_PATH ${DOMAIN}.proto
    protoc -I=$PROTO_PATH --grpc-gateway_out=paths=source_relative,grpc_api_configuration=$PROTO_PATH/${DOMAIN}.yaml:$GO_OUT_PATH ${DOMAIN}.proto
    PBTS_OUT_DIR=../coolcar2022/miniprogram/service/proto_gen/${DOMAIN}
    PBTS_BIN_DIR=../coolcar2022/miniprogram/node_modules/.bin
    $PBTS_BIN_DIR/pbjs -t static -w es6  $PROTO_PATH/${DOMAIN}.proto --no-create --no-encode --no-decode --no-verify --no-delimited -o $PBTS_OUT_DIR/${DOMAIN}_pb_tmp.js
    echo 'import * as $protobuf from "protobufjs";' > $PBTS_OUT_DIR/${DOMAIN}_pb.js
    cat $PBTS_OUT_DIR/${DOMAIN}_pb_tmp.js >> $PBTS_OUT_DIR/${DOMAIN}_pb.js
    rm $PBTS_OUT_DIR/${DOMAIN}_pb_tmp.js
    $PBTS_BIN_DIR/pbts -o $PBTS_OUT_DIR/${DOMAIN}_pb.d.ts $PBTS_OUT_DIR/${DOMAIN}_pb.js
}

genProto auth
genProto rental