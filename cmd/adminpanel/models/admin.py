from django.contrib import admin
from .models import Store, Currency, Category, Product, ProductCategory, ProductPrice, StoreSchedule

# Register your models here.

@admin.register(Store)
class StoreAdmin(admin.ModelAdmin):
    list_display = ["description", "default_currency"]
    ordering = ["description"]

@admin.register(StoreSchedule)
class StoreScheduleAdmin(admin.ModelAdmin):
    list_display = ["store", "day_of_week", "start_time", "end_time"]
    ordering = ["store"]

@admin.register(Currency)
class CurrencyAdmin(admin.ModelAdmin):
    list_display = ["code", "exp", "decimal_sep", "min_amount", "max_amount", "is_enabled"]
    search_fields = ["code"]
    ordering = ["code"]
    list_filter = ["is_enabled"]

@admin.register(Category)
class CategoryAdmin(admin.ModelAdmin):
    list_display = ["title", "is_enabled"]
    search_fields = ["title"]
    ordering = ["title"]
    list_filter = ["is_enabled"]

@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    list_display = ["title", "is_enabled", "created_at", "updated_at"]
    search_fields = ["title"]
    ordering = ["title"]
    list_filter = ["is_enabled", "created_at", "updated_at"]

@admin.register(ProductCategory)
class ProductCategoryAdmin(admin.ModelAdmin):
    list_display = ["category", "product"]
    ordering = ["category", "product"]

@admin.register(ProductPrice)
class ProductPriceAdmin(admin.ModelAdmin):
    list_display = ["product", "currency", "price"]
    ordering = ["product", "currency"]
