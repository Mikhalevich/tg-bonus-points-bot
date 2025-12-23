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

    def __str__(self):
        return self.title


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

    def __str__(self):
        return self.code


class Product(models.Model):
    title = models.TextField()
    is_enabled = models.BooleanField()
    created_at = models.DateTimeField()
    updated_at = models.DateTimeField()

    class Meta:
        managed = False
        db_table = 'product'

    def __str__(self):
        return self.title


class ProductCategory(models.Model):
    product = models.ForeignKey(Product, models.DO_NOTHING)
    category = models.ForeignKey(Category, models.DO_NOTHING)

    class Meta:
        managed = False
        db_table = 'product_category'
        unique_together = (('product', 'category'),)

    def __str__(self):
        return "{} : {}".format(self.category, self.product)


class ProductPrice(models.Model):
    product = models.ForeignKey(Product, models.DO_NOTHING)
    currency = models.ForeignKey(Currency, models.DO_NOTHING)
    price = models.IntegerField()

    class Meta:
        managed = False
        db_table = 'product_price'
        unique_together = (('product', 'currency'),)

    def __str__(self):
        return "{} : {} : {}".format(self.product.title, self.currency.code, self.price)


class Store(models.Model):
    description = models.TextField()
    default_currency = models.ForeignKey(Currency, models.DO_NOTHING)

    class Meta:
        managed = False
        db_table = 'store'

    def __str__(self):
        return self.description


class StoreSchedule(models.Model):
    class DayOfWeekEnum(models.TextChoices):
        MONDAY = 'Monday', 'Monday'
        TUESDAY = 'Tuesday', 'Tuesday'
        WEDNESDAY = 'Wednesday', 'Wednesday'
        THURSDAY = 'Thursday', 'Thursday'
        FRIDAY = 'Friday', 'Friday'
        SATURDAY = 'Saturday', 'Saturday'
        SUNDAY = 'Sunday', 'Sunday'

    store = models.ForeignKey(Store, models.DO_NOTHING)
    day_of_week = models.CharField(
        max_length=20,
        choices=DayOfWeekEnum.choices,
        default=DayOfWeekEnum.MONDAY
    )
    start_time = models.DateTimeField()
    end_time = models.DateTimeField()

    class Meta:
        managed = False
        db_table = 'store_schedule'
        unique_together = (('store', 'day_of_week'),)

    def __str__(self):
        return "{} : {}".format(self.store, self.day_of_week)
