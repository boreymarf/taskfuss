from src.classes.form_processor import FormProcessor


class TestFormProcessor:
    def test_get_value_basic(self):
        data = {"title": "Something"}
        form_processor = FormProcessor()
        form_processor.load_data(data)
        value = form_processor.get_value("title")
        assert value == "Something"

    def test_get_value_layers(self):
        data = {"first": {"second": {"third": "my_value"}}}
        form_processor = FormProcessor()
        form_processor.load_data(data)
        value = form_processor.get_value("first.second.third")
        assert value == "my_value"

    def test_get_value_multiple_string_args(self):
        data = {"first": {"second": {"third": "deep"}}}
        form_processor = FormProcessor()
        form_processor.load_data(data)
        value = form_processor.get_value("first", "second", "third")
        assert value == "deep"

    def test_get_value_with_list_index_args(self):
        data = {"items": [{"name": "apple"}, {"name": "banana"}]}
        form_processor = FormProcessor()
        form_processor.load_data(data)
        value = form_processor.get_value("items", 1, "name")
        assert value == "banana"

    def test_get_value_single_string_path_still_works(self):
        data = {"a": {"b": [10, 20, 30]}}
        form_processor = FormProcessor()
        form_processor.load_data(data)
        value = form_processor.get_value("a.b.2")
        assert value == 30

    def test_get_value_nonexistent_path(self):
        data = {"x": 1}
        form_processor = FormProcessor()
        form_processor.load_data(data)
        assert form_processor.get_value("y") is None
        assert form_processor.get_value("x", "y") is None
