import { createVuetify } from "vuetify";
import "vuetify/styles";
import AddIcon from "~icons/material-symbols/add-rounded";
import SortDescIcon from "~icons/material-symbols/arrow-downward-rounded";
import DropdownIcon from "~icons/material-symbols/arrow-drop-down-rounded";
import SortAscIcon from "~icons/material-symbols/arrow-upward-rounded";
import FileIcon from "~icons/material-symbols/attach-file-rounded";
import ClearIcon from "~icons/material-symbols/cancel-rounded";
import CheckboxOffIcon from "~icons/material-symbols/check-box-outline-blank-rounded";
import CheckboxOnIcon from "~icons/material-symbols/check-box-rounded";
import SuccessIcon from "~icons/material-symbols/check-circle-rounded";
import CheckIcon from "~icons/material-symbols/check-rounded";
import PrevIcon from "~icons/material-symbols/chevron-left-rounded";
import NextIcon from "~icons/material-symbols/chevron-right-rounded";
import CloseIcon from "~icons/material-symbols/close-rounded";
import EditIcon from "~icons/material-symbols/edit-rounded";
import ErrorIcon from "~icons/material-symbols/error-rounded";
import CheckboxIndeterminateIcon from "~icons/material-symbols/indeterminate-check-box-rounded";
import InfoIcon from "~icons/material-symbols/info-rounded";
import ExpandIcon from "~icons/material-symbols/keyboard-arrow-down-rounded";
import CollapseIcon from "~icons/material-symbols/keyboard-arrow-up-rounded";
import MenuIcon from "~icons/material-symbols/menu-rounded";
import RadioOnIcon from "~icons/material-symbols/radio-button-checked-rounded";
import RadioOffIcon from "~icons/material-symbols/radio-button-unchecked-rounded";
import MinusIcon from "~icons/material-symbols/remove-rounded";
import WarningIcon from "~icons/material-symbols/warning-rounded";

export default createVuetify({
  icons: {
    aliases: {
      collapse: CollapseIcon,
      complete: CheckIcon,
      cancel: ClearIcon,
      close: CloseIcon,
      delete: ClearIcon,
      clear: CloseIcon,
      success: SuccessIcon,
      info: InfoIcon,
      warning: WarningIcon,
      error: ErrorIcon,
      prev: PrevIcon,
      next: NextIcon,
      checkboxOn: CheckboxOnIcon,
      checkboxOff: CheckboxOffIcon,
      checkboxIndeterminate: CheckboxIndeterminateIcon,
      sortAsc: SortAscIcon,
      sortDesc: SortDescIcon,
      expand: ExpandIcon,
      menu: MenuIcon,
      subgroup: DropdownIcon,
      dropdown: DropdownIcon,
      radioOn: RadioOnIcon,
      radioOff: RadioOffIcon,
      edit: EditIcon,
      file: FileIcon,
      plus: AddIcon,
      minus: MinusIcon,
    },
  },
  display: {
    mobileBreakpoint: "md",
    thresholds: {
      xs: 0,
      sm: 600,
      md: 840,
      lg: 1145,
      xl: 1545,
      xxl: 2138,
    },
  },
  theme: {
    defaultTheme: "dark",
    utilities: false,
    themes: {
      dark: {
        dark: true,
        colors: {
          background: "#160E27",
          surface: "#251842",
          "card-background": "#332457",
          primary: "#BB86FC",
          secondary: "#FFB74D",
          accent: "#7C4DFF",
        },
      },
      light: {
        colors: {
          background: "#F5F1FA",
          "card-background": "#FFF",
          primary: "#7C4DFF",
          secondary: "#FB8C00",
          accent: "#7C4DFF",
        },
      },
    },
  },
});
