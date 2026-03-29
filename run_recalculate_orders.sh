#!/bin/bash

echo "=========================================="
echo "Recalculating Existing Orders"
echo "=========================================="
echo ""
echo "This will recalculate subtotal, tax, and total for all existing orders"
echo ""

read -p "Continue? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

echo ""
echo "Running recalculation..."
go run recalculate_existing_orders.go

if [ $? -eq 0 ]; then
    echo ""
    echo "✓ Recalculation completed successfully!"
else
    echo ""
    echo "❌ Recalculation failed!"
    exit 1
fi
