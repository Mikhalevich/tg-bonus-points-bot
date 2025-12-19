from django.contrib import admin
from .models import Store, Currency, Category, Product, ProductCategory, ProductPrice, StoreSchedule

# Register your models here.

admin.site.register(Store)
admin.site.register(StoreSchedule)

admin.site.register(Currency)

admin.site.register(Category)
admin.site.register(Product)
admin.site.register(ProductCategory)
admin.site.register(ProductPrice)
