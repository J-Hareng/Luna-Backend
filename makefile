s: 
	go run src/main.go 
up: 
	make Rprod && cd $(pC) && ng build && cd $(pS) && make Rdev && go mod vendor && git add client/* && git commit -am "Automatic commit by makefile" && git push heroku master && heroku logs --tail
bs: 
	make Rdev && cd $(pC) && echo "Building..." && ng build && echo "Build Done!" && cd $(pS) && go run src/main.go 
Rprod: 
	cp $(pC)/src/app/const_buf_prod/const.ts $(pC)/src/app && echo "pluged prod var file in"
Rdev:
	cp $(pC)/src/app/const_buf_local/const.ts $(pC)/src/app && echo "pluged Local var file in"


pC := /Users/hareng/Desktop/Luna-Frontend

pS := /home/julian/github/Luna-Backend
