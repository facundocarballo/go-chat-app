import { extendTheme } from '@chakra-ui/react';
import { colors } from './colors';
import { styles } from '.';
import { components } from './components';


const theme = extendTheme({colors, styles, components});

export default theme;
