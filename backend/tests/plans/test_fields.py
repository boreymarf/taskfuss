from src.plans.field_validation_errors import (
    IncorrectTypeError,
    ListErrors,
    MaxSizeError,
    MinSizeError,
    RequiredError,
)
from src.plans.fields import CheckboxField, ListField, StrField


class TestCheckboxField:
    def test_validation_type(self):
        field = CheckboxField()

        value = True
        errors = field.validate_value(value)
        assert not errors

        value = "True"
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], IncorrectTypeError)


class TestStrField:

    def test_required(self):
        field = StrField(required=True)

        value = "Hello World!"
        errors = field.validate_value(value)
        assert not errors

        value = ""
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], RequiredError)

    def test_validation_type(self):
        field = StrField()

        value = "Hello."
        errors = field.validate_value(value)
        assert not errors

        value = 10
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], IncorrectTypeError)

        value = ["Hello."]
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], IncorrectTypeError)

    def test_min_size(self):
        field = StrField(min_size=6)

        value = "Hello World!"
        errors = field.validate_value(value)
        assert not errors

        value = "Hello"
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], MinSizeError)

    def test_max_size(self):
        field = StrField(max_size=6)

        value = "Hello"
        errors = field.validate_value(value)
        assert not errors

        value = "Hello World!"
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], MaxSizeError)


class TestListField:

    def test_required(self):
        field = ListField(item_field=StrField(), required=True)

        value = ["Allo."]
        errors = field.validate_value(value)
        assert not errors

        value = []
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], RequiredError)

    def test_validation_type(self):
        field = ListField(item_field=StrField())

        value = ["Hello."]
        errors = field.validate_value(value)
        assert not errors

        value = "Hello."
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], IncorrectTypeError)

        value = [1]
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], ListErrors)
        assert len(errors[0].errors) == 1
        assert isinstance(errors[0].errors[0][0], IncorrectTypeError)

    def test_max_size(self):
        field = ListField(item_field=StrField(), min_size=3)

        value = ["Hello", "my", "name", "is", "Borey"]
        errors = field.validate_value(value)
        assert not errors

        value = ["Who", "cares"]
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], MinSizeError)

    def test_min_size(self):
        field = ListField(item_field=StrField(), max_size=3)

        value = ["Hello", ":3c"]
        errors = field.validate_value(value)
        assert not errors

        value = ["1", "2", "3", "4?"]
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], MaxSizeError)

    def test_internal_errors(self):
        field = ListField(item_field=StrField(min_size=3, max_size=8))

        value = ["12345"]
        errors = field.validate_value(value)
        assert not errors

        value = ["12"]
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], ListErrors)
        assert len(errors[0].errors) == 1
        assert isinstance(errors[0].errors[0][0], MinSizeError)

        value = ["12", "123456789"]
        errors = field.validate_value(value)
        assert len(errors) == 1
        assert isinstance(errors[0], ListErrors)
        assert len(errors[0].errors) == 2
        assert isinstance(errors[0].errors[0][0], MinSizeError)
        assert isinstance(errors[0].errors[1][0], MaxSizeError)
