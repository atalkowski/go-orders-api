
name="rabby"
id="15307078249801941090"
price="1999"
verbose=""
fake=""

function dft() {
  echo "$1"
}

function myrand {
  min=`dft $1 1`
  max=`dft $2 6`
  number=$(expr $min + $RANDOM$RANDOM % $max)
  echo $number
}

function getCustId() {
  case "$1" in 
  rabby) echo "C0027AB5-7777-CCCC-BEBE-123450001111";;
  andy)  echo "C003A0D1-AAAA-CCCC-BEBE-123450002222";;
  *)     echo "C001CA77-CCCC-CCCC-BEBE-C3A03A5D69ED";;
  esac
}


function getItemId() {
  case "$1" in
  1) echo '"price":1999, "item_id":"BEBEBEBE-1111-1733-1001-000090001111"';;
  2) echo '"price":586,  "item_id":"BEBEBEBE-2222-1733-1002-000090002222"';;
  3) echo '"price":1050, "item_id":"BEBEBEBE-3333-1733-1003-000090003333"';;
  *) echo '"price":1999, "item_id":"BEBEBEBE-1111-1733-1001-000090001111"';;
  esac
}


function getQuantity() {
   myrand 1 5
}


function do_get() {
  $fake curl $verbose http://localhost:3001/orders/$id
}

function do_list() {
  $fake curl $verbose http://localhost:3001/orders
}

function do_create() {
  data="{
     \"customer_id\":\"$(getCustId $name)\",
     \"line_items\": [
        {$(getItemId 1), \"quantity\":$(getQuantity)},
        {$(getItemId 2), \"quantity\":$(getQuantity)}
     ]
   }"
  echo "Using this data:"
  echo "$data"
  $fake curl $verbose -X POST http://localhost:3001/orders -d "$data"
}

function blurb() {
  echo "Use $0 to run tests on the orders API
Usage:
  $0 [options] list .... list customer order items
  $0 [options] create .. create a new order
  $0 getbyid ID ........ get the order with given ID
  

and options are:
  -v ................... use verbose mode in the curl requests
  -name xxxx ........... set name which specifies a customer Id.
                         and xxxx is either andy or rabby
  -fake ................ just show the curl request but don't execute
$*"
}


cmd="blurb"

while [ $# != 0 ]; do
  arg="$1"
  shift 
  case "$arg" in
  -v | -vv) verbose="$arg";;
  -name) name="$1"; shift;;
  -fake) fake="echo";;
  list | create) cmd="do_$arg";;
  delete | get) id="$1"; shift
     if [ "$id" = "" ]; then
       blurb "Missing ID for $arg request"
       exit 1
     fi
     cmd="do_$arg";;
  *) blurb "Don't understand $arg";;
  esac
done 

$cmd
