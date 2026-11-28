<script setup lang="ts">
import {
    ComboboxAnchor,
    ComboboxCancel,
    ComboboxContent,
    ComboboxEmpty,
    ComboboxGroup,
    ComboboxInput,
    ComboboxItem,
    ComboboxItemIndicator,
    ComboboxLabel,
    ComboboxPortal,
    ComboboxRoot,
    ComboboxSeparator,
    ComboboxTrigger,
    ComboboxViewport,
} from 'reka-ui'

// Пример данных
const items = [
    { value: 'todo', label: 'Задача' },
    { value: 'in-progress', label: 'В процессе' },
    { value: 'done', label: 'Готово' },
]
</script>

<template>
    <ComboboxRoot>
        <ComboboxAnchor class="combobox-anchor">
            <ComboboxInput class="combobox-input" />
            <ComboboxTrigger class="combobox-trigger" />
            <ComboboxCancel class="combobox-cancel" />
        </ComboboxAnchor>

        <ComboboxPortal>
            <ComboboxContent class="bg-red combobox-content z-50" position="popper" positionStrategy="fixed" side="bottom"
                align="start" :side-offset="4">
                <ComboboxViewport class="combobox-viewport">
                    <ComboboxEmpty class="combobox-empty">Ничего не найдено</ComboboxEmpty>

                    <ComboboxItem v-for="item in items" :key="item.value" :value="item.value" class="combobox-item">
                        <ComboboxItemIndicator class="combobox-item-indicator" />
                        {{ item.label }}
                    </ComboboxItem>

                    <ComboboxSeparator class="combobox-separator" />

                    <ComboboxGroup class="combobox-group">
                        <ComboboxLabel class="combobox-label">Другое</ComboboxLabel>
                        <ComboboxItem class="combobox-item" :value="'other'">
                            <ComboboxItemIndicator class="combobox-item-indicator" />
                            Другое
                        </ComboboxItem>
                    </ComboboxGroup>
                </ComboboxViewport>
            </ComboboxContent>
        </ComboboxPortal>
    </ComboboxRoot>
</template>

<style lang="scss"></style>