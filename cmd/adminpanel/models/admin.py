from django.contrib import admin
from .models import Store, Currency, Category, Product

# Register your models here.

admin.site.register(Store)

admin.site.register(Currency)

admin.site.register(Category)
admin.site.register(Product)
