# why use protobuf?
Existing serialization tools in go:
- gobs (go specific)
- ad hoc (use a seperator, like 12:1:23:2)
- json/xml (not consise, expecially for nums)

# message
message ==> struct
you can have nested messages

# tag
The " = 1", " = 2" markers on each element identify the unique "tag" that field uses in the binary encoding. 

# default val
If a field value isn't set, a default value is used: zero for numeric types, the empty string for strings, false for bools. 

# repeat
If a field is repeated, the field may be repeated any number of times (including zero).

Think of repeated fields as dynamically sized arrays.