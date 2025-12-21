from django.contrib import admin
from .models import Store, Currency, Category, Product, ProductCategory, ProductPrice, StoreSchedule

# Register your models here.

admin.site.register(Store)
admin.site.register(StoreSchedule)

@admin.register(Currency)
class CurrencyAdmin(admin.ModelAdmin):
    ordering = ["code"]

@admin.register(Category)
class CategoryAdmin(admin.ModelAdmin):
    search_fields = ["title"]
    ordering = ["title"]

@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    search_fields = ["title"]
    ordering = ["title"]

admin.site.register(ProductCategory)
admin.site.register(ProductPrice)
