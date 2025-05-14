import { toast } from 'react-toastify';

const customToast = (desc) => {
    return (
        toast(desc, {
            position: "bottom-center",
            autoClose: 1000,
            hideProgressBar: false,
            closeOnClick: false,
            pauseOnHover: true,
            draggable: true,
            progress: undefined,
            theme: "colored",
          })
    );
}

export default customToast;