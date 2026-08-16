from src.domain.field_validation_errors import (
    IncorrectTypeError,
    ListErrors,
    MaxSizeError,
    MinSizeError,
    RequiredError,
)
from src.domain.fields import CheckboxField, ListField, StrField


class TestCheckboxField:
    def test_validation_type(self):
        field = CheckboxField()
        loc = "agree"

        value = True
        errors = field.validate_value(value, loc=loc)
        assert not errors

        value = "True"
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], IncorrectTypeError)
        assert errors[0].loc == loc


class TestStrField:

    def test_required(self):
        field = StrField(required=True)
        loc = "name"

        value = "Hello World!"
        errors = field.validate_value(value, loc=loc)
        assert not errors

        value = ""
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], RequiredError)
        assert errors[0].loc == loc

    def test_validation_type(self):
        field = StrField()
        loc = "title"

        value = "Hello."
        errors = field.validate_value(value, loc=loc)
        assert not errors

        value = 10
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], IncorrectTypeError)
        assert errors[0].loc == loc

        value = ["Hello."]
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], IncorrectTypeError)
        assert errors[0].loc == loc

    def test_min_size(self):
        field = StrField(min_size=6)
        loc = "bio"

        value = "Hello World!"
        errors = field.validate_value(value, loc=loc)
        assert not errors

        value = "Hello"
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], MinSizeError)
        assert errors[0].loc == loc

    def test_max_size(self):
        field = StrField(max_size=6)
        loc = "code"

        value = "Hello"
        errors = field.validate_value(value, loc=loc)
        assert not errors

        value = "Hello World!"
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], MaxSizeError)
        assert errors[0].loc == loc


class TestListField:

    def test_required(self):
        field = ListField(item_field=StrField(), required=True)
        loc = "items"

        value = ["Allo."]
        errors = field.validate_value(value, loc=loc)
        assert not errors

        value = []
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], RequiredError)
        assert errors[0].loc == loc

    def test_validation_type(self):
        field = ListField(item_field=StrField())
        loc = "tags"

        value = ["Hello."]
        errors = field.validate_value(value, loc=loc)
        assert not errors

        value = "Hello."
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], IncorrectTypeError)
        assert errors[0].loc == loc

        value = [1]
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], ListErrors)
        assert errors[0].loc == loc
        # внутренняя ошибка: индекс 0 -> строка не прошла
        assert 0 in errors[0].errors
        inner = errors[0].errors[0]
        assert len(inner) == 1
        assert isinstance(inner[0], IncorrectTypeError)
        assert inner[0].loc == f"{loc}.0"

    def test_min_size(self):
        field = ListField(item_field=StrField(), min_size=3)
        loc = "numbers"

        value = ["Hello", "my", "name", "is", "Borey"]
        errors = field.validate_value(value, loc=loc)
        assert not errors

        value = ["Who", "cares"]
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], MinSizeError)
        assert errors[0].loc == loc

    def test_max_size(self):
        field = ListField(item_field=StrField(), max_size=3)
        loc = "numbers"

        value = ["Hello", ":3c"]
        errors = field.validate_value(value, loc=loc)
        assert not errors

        value = ["1", "2", "3", "4?"]
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], MaxSizeError)
        assert errors[0].loc == loc

    def test_internal_errors(self):
        field = ListField(item_field=StrField(min_size=3, max_size=8))
        loc = "words"

        value = ["12345"]
        errors = field.validate_value(value, loc=loc)
        assert not errors

        value = ["12"]
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], ListErrors)
        assert errors[0].loc == loc
        assert 0 in errors[0].errors
        inner = errors[0].errors[0]
        assert len(inner) == 1
        assert isinstance(inner[0], MinSizeError)
        assert inner[0].loc == f"{loc}.0"

        value = ["12", "123456789"]
        errors = field.validate_value(value, loc=loc)
        assert len(errors) == 1
        assert isinstance(errors[0], ListErrors)
        assert errors[0].loc == loc
        assert 0 in errors[0].errors and 1 in errors[0].errors
        inner0 = errors[0].errors[0]
        inner1 = errors[0].errors[1]
        assert len(inner0) == 1 and isinstance(inner0[0], MinSizeError)
        assert inner0[0].loc == f"{loc}.0"
        assert len(inner1) == 1 and isinstance(inner1[0], MaxSizeError)
        assert inner1[0].loc == f"{loc}.1"
