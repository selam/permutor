# permutor
permutor is a library to create randomly `unique` (there is no way for that but however) characters from given letters for range of start and end length


# how it's works?

when you need a `3 random letter` from `set of a-zA-Z0-9` that means you need `one element` from `set of 3 chars`, this set is a calculable and reproducible.
every generate call i select a random number between 0 and length of set (max possibility for 3 char set of AZaz09 is 238328) then i check from bitset for prevent generate same element again and again

if you give a range (lets say 3 and 6 char) after i generate all possibilities of 3 letter then i calculate length of set for 4, reset the bitset and continue to generate.

if you are gone stop your application and start again you can save state of current status into a json file and load from it.

## be a calculator
when you need a 5 char random letters for A-Za-Z0-9 (916.132.832 possibility), you need a 62**5 / 8 byte memory to (109MB) keep all possibilities in a bitset, when you need a 6 letter than you need a 6GB memory.

# Usage of permutor

```go
package main

import (
   "fmt"
   "github.com/selam/permutor"
)

func main() {
   // first parameter is a alphabet, second parameter is a min length of letters
   // third parameter is max length of letters
   rp, err := permutor.NewPermutor("abcdefgh", 3, 6)
   if err != nil {
        fmt.Printf("there is a error in parameters %s\n", err.Error());
        return
   }
   letters, posibility := rp.Generate()
   if letters == "" {
      fmt.Println("all posilible elements are already created, you can reset or finish")
   }

   fmt.Printf("%d pobisibility is %s\n", possibility, letters)

   err = rp.SaveTo("/path/to/file.json")  // so you can read from another lib writed another language
   if err != nil {
      fmt.Printf("there is a error when saving state: %s\n", err.Error())
      return // or do it something else
   }

   second, err := permutor.LoadFrom("/path/to/file.json")
   if err != nil {
          fmt.Printf("there is a error when loading state: %s\n", err.Error())
          return // or do it something else
   }

   letters, _ = second.Generate() // letters already defined
   if letters == "" {
         fmt.Println("all posilible elements are already created, you can reset or finish")
         return
   }
}
```




