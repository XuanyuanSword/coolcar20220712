db.trip.createIndex({
    "trip.accountid":1,
    "trip.status":1,
},{
    unique:true,
    partialFilterExpression:{
        "trip.status":1,
    }
})

db.account.createIndex({
    "openid":1,

},{
    unique:true,

})