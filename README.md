# Any Hasher
It could be struct hasher, but honestly, give it anything.

## Why
Because all the other hashers I found didn't quite scratch my itch, I wanted to be able to modify the structure being provided over time while still parsing original data.

## How
Steps through the properties of an object/keys of a hash, sorts them alphabetically, then packs writes everything sequentially, names and values to a given hash.Hash, only this one skips zero values so when you add a new property to a struct, then load old data with that newer struct, you still get the same hash. Only time this isn't true is if you delete or rename a property so... don't do that.

## But

### If it's skipping writing fields and values for zero values, how do you know it's safe?
Well, if you change the value to non-zero, then the name and value gets written, thus resulting in a different hash

### You could be wrong.
Quite possibly so, please feel free to prove it with a test case and I'll write a patch


## Cool feature

If you use a tag on your properties `hash:"-"` it'll skip that property no matter what's written to it, kinda useful if you're writing your hash with your data to storage.