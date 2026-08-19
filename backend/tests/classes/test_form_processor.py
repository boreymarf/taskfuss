from src.classes.form_processor import FormProcessor
from src.domain import CheckboxField
from src.domain.fields import Field, StrField, TupleField


class TestFormProcessor:

    # --------------------------------------------

    def test_get_value_basic(self):
        data = {"title": "Something"}

        form_processor = FormProcessor()
        form_processor.load_data(data)

        value = form_processor.get_value("title")
        assert value == "Something"

    def test_get_value_layers(self):
        data = {"first": {"second": "Something"}}

        form_processor = FormProcessor()
        form_processor.load_data(data)

        value = form_processor.get_value("first", "second")
        assert value == "Something"

    def test_get_value_layers_dot(self):
        data = {"first": {"second": "Something"}}

        form_processor = FormProcessor()
        form_processor.load_data(data)

        value = form_processor.get_value("first.second")
        assert value == "Something"

    def test_get_value_list(self):
        data = {"first": {"second": ["list_0", "list_1"]}}

        form_processor = FormProcessor()
        form_processor.load_data(data)

        value = form_processor.get_value("first", "second", 0)
        assert value == "list_0"

    def test_get_value_list_2(self):
        data = {"first": {"second": ["list_0", "list_1"]}}

        form_processor = FormProcessor()
        form_processor.load_data(data)

        value = form_processor.get_value("first", "second", 1)
        assert value == "list_1"

    def test_get_value_list_3(self):
        data = {"first": [[0, 1], [2, 3]]}

        form_processor = FormProcessor()
        form_processor.load_data(data)

        value = form_processor.get_value("first", 1, 0)
        assert value == 2

    # --------------------------------------------

    def test_get_field_basic(self):
        fields: dict[str, Field] = {"title": StrField()}

        form_processor = FormProcessor()
        form_processor.load_fields(fields)

        value = form_processor.get_field("title")
        assert isinstance(value, StrField)

    def test_get_field_tuple(self):
        fields: dict[str, Field] = {
            "first": TupleField(fields=(StrField(), CheckboxField()))
        }

        form_processor = FormProcessor()
        form_processor.load_fields(fields)

        value = form_processor.get_field("first", 1)
        assert isinstance(value, CheckboxField)


    def test_get_field_tuple_dot(self):
        fields: dict[str, Field] = {
            "first": TupleField(fields=(StrField(), CheckboxField()))
        }

        form_processor = FormProcessor()
        form_processor.load_fields(fields)

        value = form_processor.get_field("first.1")
        assert isinstance(value, CheckboxField)

    # --------------------------------------------

    def test_insert_value_basic(self):
        form_processor = FormProcessor()

        form_processor.insert_value("value", "title")

        assert form_processor.get_value("title") == "value"

    def test_insert_value_creates_nested_dict(self):
        form_processor = FormProcessor()

        form_processor.insert_value("value", "first", "second")

        assert form_processor.get_value("first", "second") == "value"

    def test_insert_value_with_dot_path(self):
        form_processor = FormProcessor()

        form_processor.insert_value("value", "first.second")

        assert form_processor.get_value("first.second") == "value"

    def test_insert_value_creates_list(self):
        form_processor = FormProcessor()

        form_processor.insert_value("value", "first", 0)

        assert form_processor.get_value("first", 0) == "value"
        assert isinstance(form_processor.get_value("first"), list)

    def test_insert_value_creates_nested_list(self):
        form_processor = FormProcessor()

        form_processor.insert_value("value", "first", 1, 0)

        assert form_processor.get_value("first", 1, 0) == "value"
        assert form_processor.get_value("first", 0) is None
        assert isinstance(form_processor.get_value("first", 1), list)

    def test_insert_value_creates_intermediate_none(self):
        form_processor = FormProcessor()

        form_processor.insert_value("value", "first", 2)

        assert form_processor.get_value("first", 0) is None
        assert form_processor.get_value("first", 1) is None
        assert form_processor.get_value("first", 2) == "value"

    def test_insert_value_overwrites_existing(self):
        form_processor = FormProcessor()
        form_processor.load_data({"first": {"second": "old"}})

        form_processor.insert_value("new", "first", "second")

        assert form_processor.get_value("first", "second") == "new"

    def test_insert_value_mixed_types(self):
        form_processor = FormProcessor()

        form_processor.insert_value("value", "first", 0, "name")

        assert form_processor.get_value("first", 0, "name") == "value"
        assert isinstance(form_processor.get_value("first", 0), dict)
