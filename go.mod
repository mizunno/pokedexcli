module github.com/mizunno/pokedexcli

require internal/pokeapihandler v0.0.1
replace internal/pokeapihandler => ./internal/pokeapihandler

require internal/pokecache v0.0.1
replace internal/pokecache => ./internal/pokecache

go 1.21.4
