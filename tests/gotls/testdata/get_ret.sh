#!/bin/bash
OUTPUT_PATH=$1
FUNCTION="crypto/tls.(\*Conn).Read$"


# Use readelf to get function details
# readelf -s (display the symbol table), -W(Allow output width to exceed 80 characters )
FUNCTION_DETAILS=$(readelf -sW $OUTPUT_PATH | grep $FUNCTION)
# this data look like this: 
# `echo $FUNCTION_DETALIS`
#   4335: 000000000018c2e0   928 FUNC    GLOBAL DEFAULT    1 crypto/tls.(*Conn).Read
#         offset(base 16)   size(base 10) 

# about this `awk`
# echo $FUNCTION_DETAILS | awk '{for(i=1; i<=NF; i++) print $i}'
#  4335:
#  000000000018c2e0
#  928
#  FUNC
#  GLOBAL
#  DEFAULT
#  1
#  crypto/tls.(*Conn).Read
FUNCTION_START=$(echo $FUNCTION_DETAILS | awk '{print $2}')
FUNCTION_SIZE=$(echo $FUNCTION_DETAILS | awk '{print $3}')

# Convert the hexadecimal values to decimal
FUNCTION_START_DEC=$(printf "%d" 0x$FUNCTION_START)
FUNCTION_SIZE_DEC=$(printf "%d" $FUNCTION_SIZE)

# Calculate the end address
FUNCTION_END_DEC=$((FUNCTION_START_DEC + FUNCTION_SIZE_DEC))

# Convert the decimal end address back to hexadecimal
FUNCTION_END=$(printf "%x" $FUNCTION_END_DEC)

# echo "Function Start: $FUNCTION_START"
# echo "Function End: $FUNCTION_END"


# Disassemble the function
DISASM=$(objdump -d -j .text --start-address=0x$FUNCTION_START --stop-address=0x$FUNCTION_END $OUTPUT_PATH)

# Get the start address and file offset of the .text section
SECTION_START=$(readelf -W -S $OUTPUT_PATH | awk '/.text/ {print "0x"$5}')
SECTION_OFFSET=$(readelf -W -S $OUTPUT_PATH | awk '/.text/ {print "0x"$6}')

# echo "Section Start: $SECTION_START" 
# echo "Section Offset: $SECTION_OFFSET" 


# Get the virtual addresses of the "ret" instructions
RET_ADDRESSES=$(echo "$DISASM" | grep "\sret" | awk '{print $1}' | sed 's/://g')

# Convert the virtual addresses to file offsets
for ADDR in $RET_ADDRESSES; do
    # Convert the address to decimal
    ADDR_DEC=$(printf "%d" 0x$ADDR)

    # Convert the start address and file offset to decimal
    SECTION_START_DEC=$(printf "%d" $SECTION_START)
    SECTION_OFFSET_DEC=$(printf "%d" $SECTION_OFFSET)

    # Calculate the file offset of the address
    OFFSET_DEC=$((ADDR_DEC - SECTION_START_DEC + SECTION_OFFSET_DEC))

    # Convert the file offset to hexadecimal
    OFFSET=$(printf "%x" $OFFSET_DEC)

    echo -n "0x$OFFSET, "
done