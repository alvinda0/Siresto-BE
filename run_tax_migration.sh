#!/bin/bash

echo "Running tax fields migration for orders table..."
go run add_tax_fields_to_orders.go

if [ $? -eq 0 ]; then
    echo "Migration completed successfully!"
else
    echo "Migration failed!"
    exit 1
fi
