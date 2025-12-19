# This is an auto-generated Django model module.
# You'll have to do the following manually to clean this up:
#   * Rearrange models' order
#   * Make sure each model has one field with primary_key=True
#   * Make sure each ForeignKey and OneToOneField has `on_delete` set to the desired behavior
#   * Remove `managed = False` lines if you wish to allow Django to create, modify, and delete the table
# Feel free to rename the models, but don't rename db_table values or field names.
from django.db import models


class Category(models.Model):
    title = models.TextField()
    is_enabled = models.BooleanField()

    class Meta:
        managed = False
        db_table = 'category'


class Currency(models.Model):
    code = models.TextField(unique=True)
    exp = models.IntegerField()
    decimal_sep = models.TextField()
    min_amount = models.IntegerField()
    max_amount = models.IntegerField()
    is_enabled = models.BooleanField()

    class Meta:
        managed = False
        db_table = 'currency'


class Product(models.Model):
    title = models.TextField()
    is_enabled = models.BooleanField()
    created_at = models.DateTimeField()
    updated_at = models.DateTimeField()

    class Meta:
        managed = False
        db_table = 'product'


class ProductCategory(models.Model):
    product = models.ForeignKey(Product, models.DO_NOTHING)
    category = models.ForeignKey(Category, models.DO_NOTHING)

    class Meta:
        managed = False
        db_table = 'product_category'
        unique_together = (('product', 'category'),)


class ProductPrice(models.Model):
    product = models.ForeignKey(Product, models.DO_NOTHING)
    currency_id = models.IntegerField()
    price = models.IntegerField()

    class Meta:
        managed = False
        db_table = 'product_price'
        unique_together = (('product', 'currency_id'),)


class Store(models.Model):
    description = models.TextField()
    default_currency = models.ForeignKey(Currency, models.DO_NOTHING)

    class Meta:
        managed = False
        db_table = 'store'


class StoreSchedule(models.Model):
    store = models.ForeignKey(Store, models.DO_NOTHING)
    day_of_week = models.TextField()  # This field type is a guess.
    start_time = models.DateTimeField()
    end_time = models.DateTimeField()

    class Meta:
        managed = False
        db_table = 'store_schedule'
        unique_together = (('store', 'day_of_week'),)
